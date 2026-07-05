package iotdb

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/apache/iotdb-client-go/client"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type IotDb struct {
	pool *client.SessionPool
	mtx  sync.Mutex
}

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

func NewIotDb() *IotDb {
	return &IotDb{}
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

	path, err := iotDBMetricPath(meta)
	if err != nil {
		return err
	}

	statement, err := q.session.ExecuteRawDataQuery([]string{path}, q.start, q.end)
	if err != nil {
		return err
	}
	defer statement.Close()

	series := IotDBMetricSeries{
		Name:       path,
		Properties: copyMetricProperties(meta),
	}
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		value, err := iotDBColumnValue(statement, 0)
		if err != nil {
			return err
		}
		series.Samples = append(series.Samples, IotDBMetricPoint{
			Timestamp: statement.GetTimestamp(),
			Value:     value,
		})
	}
	if err != nil {
		return err
	}

	return result.AddSeries(series)
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
	properties := meta.GetProperties()
	if len(properties) == 0 {
		return nil
	}

	copied := make(map[string]string, len(properties))
	for key, value := range properties {
		copied[key] = value
	}
	return copied
}
