package iotdb

import (
	"errors"
	"fmt"
	"nova-factory-server/app/constant/datasource"
	"nova-factory-server/app/utils/uuid"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/apache/iotdb-client-go/client"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var iotDbOnce sync.Once
var iotDbInstance *IotDb

type IotDb struct {
	pool *client.SessionPool
	mtx  sync.Mutex
}

const (
	// IotDBQuerySQLProperty 用于在通用 Query 接口中传入 IoTDB SQL。
	IotDBQuerySQLProperty = "__iotdb_sql"
	// IotDBQueryColumnProperty 用于指定 SQL 查询结果中的值列。
	IotDBQueryColumnProperty = "__iotdb_column"
)

// IotDBMetricSeries 是 IoTDB 查询返回给 MetricResult 的序列结构。
type IotDBMetricSeries struct {
	Name       string
	Properties map[string]string
	Samples    []IotDBMetricPoint
}

// IotDBMetricPoint 表示 IoTDB 查询结果中的单个时间点。
type IotDBMetricPoint struct {
	Timestamp int64
	Value     float64
}

type metricSample struct {
	kind       string
	name       string
	properties map[string]string
	ts         int64
	val        float64
}

// NewMetricSample 创建可被通用 Appender 写入的时序样本。
func NewMetricSample(kind string, properties map[string]string, timestamp int64, value float64) MetricSample {
	return metricSample{
		kind:       kind,
		properties: copyStringMap(properties),
		ts:         timestamp,
		val:        value,
	}
}

func (m metricSample) GetKind() string {
	return m.kind
}

func (m metricSample) GetProperties() map[string]string {
	return copyStringMap(m.properties)
}

func (m metricSample) GetName() string {
	if len(m.name) != 0 {
		return m.name
	}
	if len(m.properties) == 0 {
		return m.kind
	}

	keys := make([]string, 0, len(m.properties))
	for key := range m.properties {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	builder.WriteString(m.kind)
	for _, key := range keys {
		builder.WriteString("|")
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(m.properties[key])
	}
	metricBuildStr := builder.String()
	str := uuid.MakeMd5([]byte(metricBuildStr))
	m.name = fmt.Sprintf("%s_%s", datasource.MetricPrefix, str)
	return m.name
}

func (m metricSample) timestamp() int64 {
	return m.ts
}

func (m metricSample) value() float64 {
	return m.val
}

func newIotDb() *IotDb {
	return &IotDb{}
}

func GetIotDb() *IotDb {
	if iotDbInstance != nil {
		return iotDbInstance
	}
	iotDbOnce.Do(func() {
		iotDbInstance = newIotDb()
	})
	return iotDbInstance
}

// connect 连接数据库
func (i *IotDb) connect() *IotDb {
	type IotDbConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		UserName string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	var d IotDbConfig
	if err := viper.UnmarshalKey("iotdb", &d); err != nil {
		panic(err)
	}
	config := &client.PoolConfig{
		Host:     d.Host,
		Port:     d.Port,
		UserName: d.UserName,
		Password: d.Password,
	}

	pool := client.NewSessionPool(config, 3, 60000, 60000, false)
	//defer sessionPool.Close()
	i.pool = &pool
	var session client.Session
	var err error
	for {
		session, err = pool.GetSession()
		if err == nil {
			pool.PutBack(session)
			//defer pool.PutBack(session)
			return i
		}
		if err.Error() == "get session timeout" {
			zap.L().Error("get session error", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		if err != nil {
			zap.L().Error("get session error", zap.Error(err))
			i.pool.PutBack(session)
			if err.Error() == "sessionPool has closed" {
				return i
			}
			time.Sleep(5 * time.Second)
			continue
		}
	}

	return i
}

func (i *IotDb) GetSession() (client.Session, error) {
	if i == nil {
		return client.Session{}, errors.New("<UNK>")
	}
	if i.pool == nil {
		i.mtx.Lock()
		if i.pool == nil {
			i.connect()
		}
		i.mtx.Unlock()
	}
	session, err := i.pool.GetSession()
	if err == nil {
		return session, nil
	}

	if err.Error() == "get session timeout" {
		zap.L().Error("get session error", zap.Error(err))
		return client.Session{}, errors.New("sessionPool has closed")
	}
	if err != nil {
		zap.L().Error("get session error", zap.Error(err))
		i.pool.PutBack(session)
		zap.L().Error("get session", zap.Error(err))
		if err.Error() == "sessionPool has closed" {
			return session, nil

		}
		return client.Session{}, err
	}
	return session, nil
}

func (i *IotDb) PutSession(session client.Session) {
	if i == nil || i.pool == nil {
		return
	}
	i.pool.PutBack(session)
}

// Appender 将 IoTDB 适配为通用时序写入接口。
func (i *IotDb) Appender() Appender {
	return &iotDBAppender{db: i}
}

// Querier 将 IoTDB 适配为通用时序查询接口。
func (i *IotDb) Querier(startTime, endTime time.Time) (Querier, error) {
	if i == nil {
		return nil, errTSDBStorageClosed
	}

	session, err := i.GetSession()
	if err != nil {
		return nil, err
	}

	return &iotDBQuerier{
		db:      i,
		session: session,
		start:   timeToMillis(startTime),
		end:     timeToMillis(endTime),
	}, nil
}

// Close 关闭 IoTDB session pool。
func (i *IotDb) Close() error {
	if i == nil {
		return nil
	}

	i.mtx.Lock()
	defer i.mtx.Unlock()
	if i.pool != nil {
		i.pool.Close()
		i.pool = nil
	}
	return nil
}

type iotDBAppender struct {
	db      *IotDb
	samples []MetricSample
	err     error
}

func (a *iotDBAppender) Append(samples []MetricSample) error {
	if a.err != nil {
		return a.err
	}
	if a.db == nil {
		a.err = errTSDBStorageClosed
		return a.err
	}

	for _, sample := range samples {
		if _, _, err := splitIotDBMetricPath(sample); err != nil {
			a.err = err
			return err
		}
		a.samples = append(a.samples, sample)
	}
	return nil
}

func (a *iotDBAppender) Commit() error {
	if a.err != nil {
		return a.err
	}
	if len(a.samples) == 0 {
		return nil
	}

	session, err := a.db.GetSession()
	if err != nil {
		return err
	}
	defer a.db.PutSession(session)

	deviceIDs := make([]string, 0, len(a.samples))
	measurements := make([][]string, 0, len(a.samples))
	dataTypes := make([][]client.TSDataType, 0, len(a.samples))
	values := make([][]interface{}, 0, len(a.samples))
	timestamps := make([]int64, 0, len(a.samples))

	for _, sample := range a.samples {
		deviceID, measurement, err := splitIotDBMetricPath(sample)
		if err != nil {
			return err
		}

		deviceIDs = append(deviceIDs, deviceID)
		measurements = append(measurements, []string{measurement})
		dataTypes = append(dataTypes, []client.TSDataType{client.DOUBLE})
		values = append(values, []interface{}{sample.value()})
		timestamps = append(timestamps, sample.timestamp())
	}

	status, err := session.InsertRecords(deviceIDs, measurements, dataTypes, values, timestamps)
	if err != nil {
		return err
	}
	return client.VerifySuccess(status)
}

type iotDBQuerier struct {
	db      *IotDb
	session client.Session
	start   int64
	end     int64
}

func (q *iotDBQuerier) Query(meta MetricMeta, hints *QueryHints, result MetricResult) error {
	if q == nil || q.db == nil {
		return errTSDBStorageClosed
	}
	if result == nil {
		return fmt.Errorf("%w: nil metric result", errTSDBInvalidMetric)
	}
	if sql := meta.GetProperties()[IotDBQuerySQLProperty]; sql != "" {
		return q.querySQL(meta, sql, meta.GetProperties()[IotDBQueryColumnProperty], result)
	}

	path, err := iotDBMetricPath(meta)
	if err != nil {
		return err
	}

	statement, err := q.session.ExecuteRawDataQuery([]string{path}, q.start, q.end)
	if err != nil {
		return err
	}
	defer statement.Close()

	return addStatementSeries(statement, path, copyMetricProperties(meta), "", result)
}

func (q *iotDBQuerier) querySQL(meta MetricMeta, sql string, column string, result MetricResult) error {
	var timeout int64 = 5000
	statement, err := q.session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		return err
	}
	defer statement.Close()

	return addStatementSeries(statement, meta.GetKind(), copyMetricProperties(meta), column, result)
}

func (q *iotDBQuerier) QueryAndClose(meta MetricMeta, hints *QueryHints, result MetricResult) error {
	defer q.Close()
	return q.Query(meta, hints, result)
}

func (q *iotDBQuerier) Close() {
	if q == nil || q.db == nil {
		return
	}
	q.db.PutSession(q.session)
	q.db = nil
}

func splitIotDBMetricPath(meta MetricMeta) (string, string, error) {
	path, err := iotDBMetricPath(meta)
	if err != nil {
		return "", "", err
	}

	index := strings.LastIndex(path, ".")
	if index <= 0 || index == len(path)-1 {
		return "", "", fmt.Errorf("%w: invalid iotdb path %q", errTSDBInvalidMetric, path)
	}
	return path[:index], path[index+1:], nil
}

func iotDBMetricPath(meta MetricMeta) (string, error) {
	if meta == nil {
		return "", fmt.Errorf("%w: nil metric meta", errTSDBInvalidMetric)
	}

	path := strings.TrimSpace(meta.GetKind())
	if path == "" {
		return "", fmt.Errorf("%w: empty iotdb path", errTSDBInvalidMetric)
	}
	return path, nil
}

func addStatementSeries(statement *client.SessionDataSet, name string, properties map[string]string, column string, result MetricResult) error {
	series := IotDBMetricSeries{
		Name:       name,
		Properties: properties,
	}
	for next, err := statement.Next(); ; next, err = statement.Next() {
		if err != nil {
			return err
		}
		if !next {
			break
		}
		value, err := iotDBColumnValueByName(statement, column)
		if err != nil {
			return err
		}
		series.Samples = append(series.Samples, IotDBMetricPoint{
			Timestamp: statement.GetTimestamp(),
			Value:     value,
		})
	}
	return result.AddSeries(series)
}

func iotDBColumnValueByName(statement *client.SessionDataSet, column string) (float64, error) {
	columnIndex := 0
	if column != "" {
		columnIndex = -1
		for index, columnName := range statement.GetColumnNames() {
			if columnName == column {
				columnIndex = index
				break
			}
		}
		if columnIndex < 0 {
			return 0, fmt.Errorf("%w: column %q not found", errTSDBInvalidMetric, column)
		}
	}
	return iotDBColumnValue(statement, columnIndex)
}

func iotDBColumnValue(statement *client.SessionDataSet, columnIndex int) (float64, error) {
	columnName := statement.GetColumnName(columnIndex)
	if statement.IsNull(columnName) {
		return 0, nil
	}

	switch statement.GetColumnDataType(columnIndex) {
	case client.BOOLEAN:
		if statement.GetBool(columnName) {
			return 1, nil
		}
		return 0, nil
	case client.INT32:
		return float64(statement.GetInt32(columnName)), nil
	case client.INT64:
		return float64(statement.GetInt64(columnName)), nil
	case client.FLOAT:
		return float64(statement.GetFloat(columnName)), nil
	case client.DOUBLE:
		return statement.GetDouble(columnName), nil
	default:
		return 0, fmt.Errorf("%w: unsupported iotdb data type %v", errTSDBInvalidMetric, statement.GetColumnDataType(columnIndex))
	}
}

func copyMetricProperties(meta MetricMeta) map[string]string {
	return copyStringMap(meta.GetProperties())
}

func copyStringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return nil
	}

	copied := make(map[string]string, len(values))
	for key, value := range values {
		copied[key] = value
	}
	return copied
}
