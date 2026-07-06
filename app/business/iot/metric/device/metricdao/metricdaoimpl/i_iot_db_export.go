package metricdaoimpl

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/iotdb-client-go/client"
	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"go.uber.org/zap"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/math"
	"nova-factory-server/app/utils/time"
	"nova-factory-server/app/utils/uuid"
	"sort"
	"strings"
	stdtime "time"
)

type iotDbExport struct {
	iotDb *iotdb.IotDb
}

func newIotDbExport(iotDb *iotdb.IotDb) iDaoExport {
	i := &iotDbExport{
		iotDb: iotDb,
	}
	i.init()
	return i
}

func (i *iotDbExport) init() {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("iotdb.GetSession()", zap.Error(err))
		panic(err)
	}
	defer i.iotDb.PutSession(session)
	for {
		statement, err := session.ExecuteStatement("count databases root.device")
		if err != nil {
			zap.L().Error("execute statement", zap.Error(err))
			stdtime.Sleep(1 * stdtime.Second)
			continue
		}
		hasDatabase, err := statement.Next()
		if err != nil {
			zap.L().Error("get hasDatabase", zap.Error(err))
			stdtime.Sleep(1 * stdtime.Second)
			continue
		}

		count := statement.GetInt32("count")
		if count < 1 {
			session.ExecuteStatement("create database root.device")
			stdtime.Sleep(1 * stdtime.Second)
			continue
		}

		statement, err = session.ExecuteStatement("count databases root.run_status_device")
		if err != nil {
			stdtime.Sleep(1 * stdtime.Second)
			return
		}
		hasDatabase, err = statement.Next()
		if err != nil {
			zap.L().Error("get hasDatabase", zap.Error(err))
			stdtime.Sleep(1 * stdtime.Second)
			continue
		}

		if !hasDatabase {
			stdtime.Sleep(1 * stdtime.Second)
			continue
		}

		count = statement.GetInt32("count")
		if count < 1 {
			session.ExecuteStatement("create database root.run_status_device")
			stdtime.Sleep(1 * stdtime.Second)
			continue
		}

		break
	}

	// 创建设备数据采集模板
	session.ExecuteStatement(fmt.Sprintf("create device template %s ALIGNED (value DOUBLE)", iotdb2.NOVA_DEVICE_TEMPLATE))
	// 创建设备运行时间统计模板
	session.ExecuteStatement(fmt.Sprintf("create device template %s ALIGNED (duration INT64, status INT64)", iotdb2.NOVA_DEVICE_RUN_TEMPLATE))

}

type iotMetricMeta struct {
	kind       string
	properties map[string]string
	name       string
}

func (m iotMetricMeta) GetKind() string {
	return m.kind
}

func (m iotMetricMeta) GetProperties() map[string]string {
	return m.properties
}

func (m iotMetricMeta) GetName() string {
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
	return str
}

type iotMetricQueryResult struct {
	data *metricmodels.MetricQueryData
}

func (r *iotMetricQueryResult) GetName() string {
	//TODO implement me
	return ""
}

func (r *iotMetricQueryResult) GetKind() string {
	return ""
}

func (r *iotMetricQueryResult) GetProperties() map[string]string {
	return nil
}

func (r *iotMetricQueryResult) AddSeries(series iotdb.MetricSeries) error {
	iotSeries, ok := series.(iotdb.IotDBMetricSeries)
	if !ok {
		return fmt.Errorf("unsupported iotdb metric series %T", series)
	}
	for _, sample := range iotSeries.Samples {
		r.data.Values = append(r.data.Values, metricmodels.MetricQueryValue{
			Time:  sample.Timestamp,
			Value: math.RoundFloat(sample.Value, 2),
		})
	}
	return nil
}

func (i *iotDbExport) Export(ctx context.Context, data []*metricmodels.NovaMetricsDevice) error {
	if len(data) == 0 {
		return nil
	}

	samples := make([]iotdb.MetricSample, 0, len(data))
	for _, value := range data {
		if value == nil || value.StartTimeUnix == nil {
			continue
		}
		name := iotdb2.MakeDeviceTemplateName(int64(value.DeviceId), int64(value.TemplateId), int64(value.DataId)) + ".value"
		samples = append(samples, iotdb.NewMetricSample(name, nil, value.StartTimeUnix.UnixMilli(), value.Value))
	}
	if len(samples) == 0 {
		return nil
	}

	appender := i.iotDb.Appender()
	if err := appender.Append(samples); err != nil {
		zap.L().Error("iotdb appender append error", zap.Error(err))
		return err
	}
	if err := appender.Commit(); err != nil {
		zap.L().Error("iotdb appender commit error", zap.Error(err))
		return err
	}
	return nil
}

func (i *iotDbExport) Metric(c *gin.Context, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
	if req == nil {
		return nil, nil
	}

	var startTime string
	if req.Start > 0 {
		startTime = time.GetStartTime(req.Start, 200)
	}
	endTime := time.GetEndTimeUseNow(req.End, true)

	if startTime == "" {
		return nil, errors.New("开始时间不能为空")
	}

	if endTime == "" {
		return nil, errors.New("结束时间不能为空")
	}
	if req.Step <= 0 {
		req.Step = 1
	}

	name := iotdb2.MakeDeviceTemplateName(int64(req.DeviceId), int64(req.TemplateId), int64(req.DataId))
	data := metricmodels.NewMetricQueryData()
	data.Id = name
	sql := fmt.Sprintf("select avg(value) as value from %s group by([%s, %s), %dm, %dm);",
		name, startTime, endTime, req.Step, req.Step)

	querier, err := i.iotDb.Querier(stdtime.UnixMilli(int64(req.Start)), metricEndTime(req.End))
	if err != nil {
		return nil, err
	}
	result := &iotMetricQueryResult{data: data}
	meta := iotMetricMeta{
		kind: name,
		properties: map[string]string{
			iotdb.IotDBQuerySQLProperty:    sql,
			iotdb.IotDBQueryColumnProperty: "value",
		},
	}
	if err := querier.QueryAndClose(meta, nil, result); err != nil {
		zap.L().Error("iotdb querier query error", zap.Error(err))
		return nil, err
	}
	return data, nil
}

// InstallDevice 安装设备模板
func (i *iotDbExport) InstallDevice(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := iotdb2.MakeDeviceTemplateName(deviceId, device.TemplateID, device.DeviceConfigID)
	// 创建设备模板
	group, err := session.SetStorageGroup(name)
	if err != nil {
		zap.L().Error("创建设备数据库失败, ", zap.Error(err), zap.Any("code", group.GetCode()))
		return err
	}

	// 挂载设备模板
	_, err = session.ExecuteStatement(fmt.Sprintf("set device template %s to %s", iotdb2.NOVA_DEVICE_TEMPLATE, name))
	if err != nil {
		zap.L().Error("绑定设备数据库失败, ", zap.Error(err))
		return err
	}

	// 激活设备模板
	_, err = session.ExecuteStatement(fmt.Sprintf("create timeseries using device template on %s", name))
	if err != nil {
		zap.L().Error("激活设备模板失败, ", zap.Error(err))
		return err
	}
	return nil
}

// InstallRunStatusDevice 运行状态设备模板
func (i *iotDbExport) InstallRunStatusDevice(c *gin.Context, deviceId int64) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := iotdb2.MakeRunDeviceTemplateName(deviceId)
	// 创建设备模板
	group, err := session.SetStorageGroup(name)
	if err != nil {
		zap.L().Error("创建设备数据库失败, ", zap.Error(err), zap.Any("code", group.GetCode()))
		return err
	}

	// 挂载设备模板
	_, err = session.ExecuteStatement(fmt.Sprintf("set device template %s to %s", iotdb2.NOVA_DEVICE_RUN_TEMPLATE, name))
	if err != nil {
		zap.L().Error("绑定设备数据库失败, ", zap.Error(err))
		return err
	}

	// 激活设备模板
	_, err = session.ExecuteStatement(fmt.Sprintf("create timeseries using device template on %s", name))
	if err != nil {
		zap.L().Error("激活设备模板失败, ", zap.Error(err))
		return err
	}
	return nil
}

// UnInStallRunStatusDevice 卸载设备运行状态模板
func (i *iotDbExport) UnInStallRunStatusDevice(c *gin.Context, deviceId int64) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := fmt.Sprintf(iotdb2.ROOT_RUN_STATUS_DEVICE_TEMPLATE_NAME, deviceId)

	// 删除模板表示的某一组时间序列
	_, err = session.ExecuteStatement(fmt.Sprintf("deactivate device template %s from %s", iotdb2.NOVA_DEVICE_RUN_TEMPLATE, name))
	if err != nil {
		zap.L().Error("deactivate  device template", zap.Error(err))
		return err
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("unset device template %s from %s", iotdb2.NOVA_DEVICE_RUN_TEMPLATE, name))
	if err != nil {
		zap.L().Error("unset  device template", zap.Error(err))
		return err
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("drop database %s", name))
	if err != nil {
		zap.L().Error("unset  device template", zap.Error(err))
		return err
	}
	return nil
}

// UnInStallDevice 卸载设备模板
func (i *iotDbExport) UnInStallDevice(c *gin.Context, deviceId int64, templateId int64, dataId int64) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := iotdb2.MakeDeviceTemplateName(deviceId, templateId, dataId)

	// 删除模板表示的某一组时间序列
	_, err = session.ExecuteStatement(fmt.Sprintf("deactivate device template %s from %s", iotdb2.NOVA_DEVICE_TEMPLATE, name))
	if err != nil {
		zap.L().Error("deactivate  device template", zap.Error(err))
		return err
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("unset device template %s from %s", iotdb2.NOVA_DEVICE_TEMPLATE, name))
	if err != nil {
		zap.L().Error("unset  device template", zap.Error(err))
		return err
	}

	_, err = session.ExecuteStatement(fmt.Sprintf("drop database %s", name))
	if err != nil {
		zap.L().Error("unset  device template", zap.Error(err))
		return err
	}
	return nil
}

// Predict 趋势预测
func (i *iotDbExport) Predict(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
	if req == nil {
		return nil, nil
	}
	// call inference(_STLForecaster, "select avg(value) from root.device.dev10366836907956274839,root.device.dev18109908223314572656 group by([2025-07-10 11:20:32, 2025-07-10 15:20:32), 3m, 3m) order by time desc  align by time", generateTime=True,predict_length=10);

	session, err := i.iotDb.GetSession()
	if err != nil {
		return nil, err
	}
	defer i.iotDb.PutSession(session)

	var startTime string
	if req.Start > 0 {
		startTime = time.GetStartTime(req.Start, 200)
	}
	endTime := time.GetEndTimeUseNow(req.End, true)

	if startTime == "" {
		return nil, errors.New("开始时间不能为空")
	}

	if endTime == "" {
		return nil, errors.New("结束时间不能为空")
	}
	if req.Step <= 0 {
		req.Step = 1
	}
	name := iotdb2.MakeDeviceTemplateName(int64(req.DeviceId), int64(req.TemplateId), int64(req.DataId))
	var timeout int64 = 5000
	var data *metricmodels.MetricQueryData = metricmodels.NewMetricQueryData()
	if device.AggFunction == "" {
		device.AggFunction = "avg"
	}
	sql := fmt.Sprintf("select %s(value) as value from %s group by([%s, %s), %dm, %dm)", device.AggFunction,
		name, startTime, endTime, req.Step, req.Step)
	predictSql := fmt.Sprintf("call inference(stl_forecaster, \"%s\", generateTime=True, predict_length=10)",
		sql)
	// select avg(value) from root.device.dev375986234780028928 group by([2025-07-07 20:52:28, 2025-07-07 21:52:28), 3m, 3m);
	statement, err := session.ExecuteQueryStatement(predictSql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		v := statement.GetDouble("output0")
		data.Values = append(data.Values, metricmodels.MetricQueryValue{
			Time:  timestamp,
			Value: math.RoundFloat(v, 2),
		})

	}
	data.Id = name
	return data, nil
}

func (i *iotDbExport) List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error) {
	var startTime string
	if req.Start > 0 {
		startTime = time.GetStartTime(req.Start, 200)
	}
	endTime := time.GetEndTimeUseNow(req.End, true)

	if startTime == "" {
		return nil, errors.New("开始时间不能为空")
	}

	if endTime == "" {
		return nil, errors.New("结束时间不能为空")
	}

	if req.Size <= 0 {
		req.Size = 20
	}

	if req.Size > 50 {
		req.Size = 50
	}

	if req.Page < 1 {
		req.Page = 1
	}

	session, err := i.iotDb.GetSession()
	if err != nil {
		return nil, err
	}
	defer i.iotDb.PutSession(session)

	offset := (req.Page - 1) * req.Size

	where := ""

	if startTime != "" {
		where = where + fmt.Sprintf(" time >= %s", startTime)
	}

	if endTime != "" {
		if where == "" {
			where = where + fmt.Sprintf(" time < %s", endTime)
		} else {
			where = where + fmt.Sprintf(" and time < %s", endTime)
		}
	}

	tableName := "root.device.*"
	if len(req.Dev) != 0 {
		tableName = strings.Join(req.Dev, ",")
	}

	if where != "" {
		where = fmt.Sprintf(" where %s", where)
	}
	var resp devicemonitormodel.DevDataResp
	resp.Rows = make([]devicemonitormodel.DevData, 0)
	var timeout int64 = 5000
	sql := fmt.Sprintf("select * from %s %s order by time desc limit %d offset %d align by device",
		tableName, where, req.Size, offset)
	// select avg(value) from root.device.dev375986234780028928 group by([2025-07-07 20:52:28, 2025-07-07 21:52:28), 3m, 3m);
	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		device := statement.GetText("Device")
		value := statement.GetDouble("value")
		resp.Rows = append(resp.Rows, devicemonitormodel.DevData{
			Time:  time.MillToTime(timestamp),
			Value: value,
			Dev:   device,
		})

	}

	return &resp, nil
}

func (i *iotDbExport) Count(c *gin.Context, req *devicemonitormodel.DevDataReq) (uint64, error) {
	var startTime string
	if req.Start > 0 {
		startTime = time.GetStartTime(req.Start, 200)
	}
	endTime := time.GetEndTimeUseNow(req.End, true)

	if startTime == "" {
		return 0, errors.New("开始时间不能为空")
	}

	if endTime == "" {
		return 0, errors.New("结束时间不能为空")
	}

	if req.Size <= 0 {
		req.Size = 20
	}

	if req.Size > 50 {
		req.Size = 50
	}

	if req.Page < 1 {
		req.Page = 1
	}

	session, err := i.iotDb.GetSession()
	if err != nil {
		return 0, err
	}
	defer i.iotDb.PutSession(session)

	where := ""

	if startTime != "" {
		where = where + fmt.Sprintf(" time >= %s", startTime)
	}

	if endTime != "" {
		if where == "" {
			where = where + fmt.Sprintf(" time < %s", endTime)
		} else {
			where = where + fmt.Sprintf(" and time < %s", endTime)
		}
	}

	tableName := "root.device.*"
	if len(req.Dev) != 0 {
		tableName = strings.Join(req.Dev, ",")
	}

	if where != "" {
		where = fmt.Sprintf(" where %s", where)
	}
	var resp devicemonitormodel.DevDataResp
	resp.Rows = make([]devicemonitormodel.DevData, 0)
	var timeout int64 = 5000

	sql := fmt.Sprintf("select count(*) from %s %s order by time desc align by device",
		tableName, where)
	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return 0, err
	}
	var sum uint64 = 0
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		count := statement.GetInt64("count(value)")
		sum = sum + uint64(count)
	}
	resp.Total = sum
	return resp.Total, nil
}

// Query dashboard 查询接口
func (i *iotDbExport) Query(c *gin.Context, req *metricmodels.MetricDataQueryReq) (*metricmodels.MetricQueryData, error) {
	if req == nil {
		return nil, nil
	}

	session, err := i.iotDb.GetSession()
	if err != nil {
		return nil, err
	}
	defer i.iotDb.PutSession(session)

	var startTime string
	if req.Start > 0 {
		startTime = time.GetStartTime(req.Start, 200)
	}
	endTime := time.GetEndTimeUseNow(req.End, true)

	if startTime == "" {
		return nil, errors.New("开始时间不能为空")
	}

	if endTime == "" {
		return nil, errors.New("结束时间不能为空")
	}

	var interval int
	if req.Interval == 0 {
		interval = int((req.End - req.Start) / 60 / 30 / 1000)
	} else {
		interval = req.Interval
	}

	if interval == 0 {
		interval = 1
	}

	if req.Field == "" {
		req.Field = " as value"
	}

	var having string
	if req.Having != "" {
		having = " having " + req.Having
	}

	var level string
	if req.Level != nil {
		levelValue := *req.Level
		level = fmt.Sprintf(",level=%d", levelValue)
	}

	var timeout int64 = 5000
	var data *metricmodels.MetricQueryData = metricmodels.NewMetricQueryData()
	var sql string
	if req.Type == "bar" || req.Type == "line" || req.Type == "area" || req.Type == "toplist" {
		if req.Step != 0 {
			sql = fmt.Sprintf("select %s %s from %s group by([%s, %s), %dm, %dm)%s %s",
				req.Expression, req.Field, req.Name, startTime, endTime, interval, req.Step, level, having)
		} else {
			sql = fmt.Sprintf("select %s %s from %s group by([%s, %s), %dm)%s %s",
				req.Expression, req.Field, req.Name, startTime, endTime, interval, level, having)
		}
	} else {
		sql = fmt.Sprintf("select %s %s from %s where time > %s and time < %s",
			req.Expression, req.Field, req.Name, startTime, endTime)
	}

	if req.Predict.Param == "" {
		req.Predict.Param = "generateTime=True, predict_length=30"
	}

	if req.Predict.Enable {
		sql = fmt.Sprintf("call inference(%s, \"%s\", %s)",
			req.Predict.Model, sql, req.Predict.Param)
	}
	zap.L().Info("iotdb search sql", zap.String("sql", sql))
	// select avg(value) from root.device.dev375986234780028928 group by([2025-07-07 20:52:28, 2025-07-07 21:52:28), 3m, 3m);
	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}

	if len(statement.GetColumnNames()) <= 1 {
		for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
			timestamp := statement.GetTimestamp()
			var v float64
			dataType := statement.GetColumnDataType(0)
			switch dataType {
			case client.BOOLEAN:
				{
					dataValue := statement.GetBool(statement.GetColumnName(0))
					if dataValue == true {
						v = 1.0
					} else {
						v = 0.0
					}
					break
				}
			case client.INT32:
				{
					dataValue := statement.GetInt32(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}
			case client.INT64:
				{
					dataValue := statement.GetInt64(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}
			case client.FLOAT:
				{
					dataValue := statement.GetFloat(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}
			case client.DOUBLE:
				{
					dataValue := statement.GetDouble(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}

			}

			data.Values = append(data.Values, metricmodels.MetricQueryValue{
				Time:  timestamp,
				Value: math.RoundFloat(v, 2),
			})

		}
	} else {
		data.MultiValues = make([][]metricmodels.MetricQueryValue, len(statement.GetColumnNames()))

		for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
			timestamp := statement.GetTimestamp()

			for k, column := range statement.GetColumnNames() {
				var v float64
				dataType := statement.GetColumnDataType(0)
				switch dataType {
				case client.BOOLEAN:
					{
						dataValue := statement.GetBool(column)
						if dataValue == true {
							v = 1.0
						} else {
							v = 0.0
						}
						break
					}
				case client.INT32:
					{
						dataValue := statement.GetInt32(column)
						v = float64(dataValue)
						break
					}
				case client.INT64:
					{
						dataValue := statement.GetInt64(column)
						v = float64(dataValue)
						break
					}
				case client.FLOAT:
					{
						dataValue := statement.GetFloat(column)
						v = float64(dataValue)
						break
					}
				case client.DOUBLE:
					{
						dataValue := statement.GetDouble(column)
						v = float64(dataValue)
						break
					}

				}
				data.MultiValues[k] = append(data.MultiValues[k], metricmodels.MetricQueryValue{
					Time:  timestamp,
					Value: math.RoundFloat(v, 2),
				})
			}
		}
	}

	return data, nil
}

func (i *iotDbExport) CounterByTimeRange(startTime int64, endTime int64, interval string) (*metricmodels.MetricQueryData, error) {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer i.iotDb.PutSession(session)

	if interval == "" {
		intervalValue := (endTime - startTime) / 60 / 30 / 1000
		interval = fmt.Sprintf("%dm", intervalValue)
	}

	var timeout int64 = 5000
	var data *metricmodels.MetricQueryData = metricmodels.NewMetricQueryData()

	sql := fmt.Sprintf("select count(value) as value from root.device.** group by([%s, %s), %s),level=1",
		time.GetStartTime(uint64(startTime), 0), time.GetEndTime(uint64(endTime), 0), interval)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}
	if len(statement.GetColumnNames()) <= 1 {
		for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
			timestamp := statement.GetTimestamp()
			var v float64
			dataType := statement.GetColumnDataType(0)
			switch dataType {
			case client.BOOLEAN:
				{
					dataValue := statement.GetBool(statement.GetColumnName(0))
					if dataValue == true {
						v = 1.0
					} else {
						v = 0.0
					}
					break
				}
			case client.INT32:
				{
					dataValue := statement.GetInt32(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}
			case client.INT64:
				{
					dataValue := statement.GetInt64(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}
			case client.FLOAT:
				{
					dataValue := statement.GetFloat(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}
			case client.DOUBLE:
				{
					dataValue := statement.GetDouble(statement.GetColumnName(0))
					v = float64(dataValue)
					break
				}

			}

			data.Values = append(data.Values, metricmodels.MetricQueryValue{
				Time:  timestamp,
				Value: math.RoundFloat(v, 2),
			})

		}
	} else {
		data.MultiValues = make([][]metricmodels.MetricQueryValue, len(statement.GetColumnNames()))

		for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
			timestamp := statement.GetTimestamp()

			for k, column := range statement.GetColumnNames() {
				var v float64
				dataType := statement.GetColumnDataType(0)
				switch dataType {
				case client.BOOLEAN:
					{
						dataValue := statement.GetBool(column)
						if dataValue == true {
							v = 1.0
						} else {
							v = 0.0
						}
						break
					}
				case client.INT32:
					{
						dataValue := statement.GetInt32(column)
						v = float64(dataValue)
						break
					}
				case client.INT64:
					{
						dataValue := statement.GetInt64(column)
						v = float64(dataValue)
						break
					}
				case client.FLOAT:
					{
						dataValue := statement.GetFloat(column)
						v = float64(dataValue)
						break
					}
				case client.DOUBLE:
					{
						dataValue := statement.GetDouble(column)
						v = float64(dataValue)
						break
					}

				}
				data.MultiValues[k] = append(data.MultiValues[k], metricmodels.MetricQueryValue{
					Time:  timestamp,
					Value: math.RoundFloat(v, 2),
				})
			}
		}
	}

	return data, nil
}

func (i *iotDbExport) CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*devicemonitormodel.TypeDeviceCounterRank, error) {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer i.iotDb.PutSession(session)

	var timeout int64 = 5000

	sql := fmt.Sprintf("select count(value) from root.device.** where time > %s and time < %s  order by count(value) desc limit %d ALIGN BY DEVICE",
		time.GetStartTime(uint64(startTime), 0), time.GetEndTime(uint64(endTime), 0), limit)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}

	rank := devicemonitormodel.TypeDeviceCounterRank{
		Rows: make([]*devicemonitormodel.TypeDeviceCounterRankValue, 0),
	}

	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		device := statement.GetText(statement.GetColumnName(0))
		value := statement.GetInt64(statement.GetColumnName(1))
		rank.Rows = append(rank.Rows, &devicemonitormodel.TypeDeviceCounterRankValue{
			Time:  timestamp,
			Dev:   device,
			Value: value,
		})
	}

	return &rank, nil
}

func (i *iotDbExport) StatDeviceStatus(c *gin.Context, startTime string, endTime string,
	status int) (*devicemonitormodel.DeviceStatusList, error) {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer i.iotDb.PutSession(session)

	var timeout int64 = 5000

	sql := fmt.Sprintf("select sum(duration) as value from root.run_status_device.** where time > %s and time < %s and status = %d align by device",
		startTime, endTime, status)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}
	data := devicemonitormodel.NewDeviceStatusList()
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		v := statement.GetDouble("value")
		deviceName := statement.GetText("Device")
		var deviceId int64
		_, err := fmt.Sscanf(deviceName, "root.run_status_device.dev%d", &deviceId)
		if err != nil {
			zap.L().Error("fmt Sscanf error", zap.Error(err))
			continue
		}
		data.List = append(data.List, devicemonitormodel.DeviceStatus{
			DeviceId: deviceId,
			Value:    v,
			Status:   status,
		})
	}

	return data, nil
}

func (i *iotDbExport) StatDeviceProcess(c *gin.Context, startTime string, endTime string, interval string,
	status int) (*devicemonitormodel.DeviceProcessList, error) {
	var processList devicemonitormodel.DeviceProcessList
	processList.List = make(map[string][]devicemonitormodel.DeviceStatus)
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer i.iotDb.PutSession(session)

	var timeout int64 = 5000

	sql := fmt.Sprintf("select sum(duration) as value from root.run_status_device.** where status = %d group by ([%s, %s), %s)  align by device",
		status, startTime, endTime, interval)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("读取设备运行过程失败:", zap.Error(err))
		return nil, err
	}
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		deviceName := statement.GetText(statement.GetColumnName(0))
		duration := statement.GetDouble(statement.GetColumnName(1))
		_, ok := processList.List[deviceName]
		if !ok {
			processList.List[deviceName] = make([]devicemonitormodel.DeviceStatus, 0)
		}
		processList.List[deviceName] = append(processList.List[deviceName], devicemonitormodel.DeviceStatus{
			Value:  duration,
			Status: status,
			Time:   timestamp,
		})
		continue
	}

	return &processList, nil
}

func (i *iotDbExport) StatDeviceRunStatus(c *gin.Context, startTime string,
	endTime string) ([]devicemonitormodel.DeviceRunStat, error) {
	var runStatList []devicemonitormodel.DeviceRunStat = make([]devicemonitormodel.DeviceRunStat, 0)
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer i.iotDb.PutSession(session)

	var timeout int64 = 5000

	sql := fmt.Sprintf("select last_value(status) from root.run_status_device.** where time>=%s and time < %s  align by device;",
		startTime, endTime)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("读取设备运行过程失败:", zap.Error(err))
		return nil, err
	}
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		deviceName := statement.GetText(statement.GetColumnName(0))
		status := statement.GetInt64(statement.GetColumnName(1))
		var stat devicemonitormodel.DeviceRunStat = devicemonitormodel.DeviceRunStat{
			Time:   timestamp,
			Status: int((status)),
			Dev:    deviceName,
		}
		runStatList = append(runStatList, stat)
		continue
	}

	return runStatList, nil
}

func (i *iotDbExport) StatDeviceStatusByDeviceId(c *gin.Context, startTime string, endTime string,
	deviceId int64, status int) (*devicemonitormodel.DeviceStatusList, error) {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer i.iotDb.PutSession(session)
	deviceKey := iotdb2.MakeRunDeviceTemplateName(deviceId)
	var timeout int64 = 5000

	sql := fmt.Sprintf("select sum(duration) as value from %s where time > %s and time < %s and status = %d align by device",
		deviceKey, startTime, endTime, status)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}
	data := devicemonitormodel.NewDeviceStatusList()
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		v := statement.GetDouble("value")
		deviceName := statement.GetText("Device")
		var deviceId int64
		_, err := fmt.Sscanf(deviceName, "root.run_status_device.dev%d", &deviceId)
		if err != nil {
			zap.L().Error("fmt Sscanf error", zap.Error(err))
			continue
		}
		data.List = append(data.List, devicemonitormodel.DeviceStatus{
			DeviceId: deviceId,
			Value:    v,
			Status:   status,
		})
	}

	return data, nil
}

func (i *iotDbExport) StatDeviceProcessByDeviceId(c *gin.Context, startTime string, endTime string,
	deviceId int64, interval string,
	status int) (*devicemonitormodel.DeviceProcessList, error) {
	var processList devicemonitormodel.DeviceProcessList
	processList.List = make(map[string][]devicemonitormodel.DeviceStatus)
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return nil, err
	}
	defer i.iotDb.PutSession(session)

	var timeout int64 = 5000

	deviceKey := iotdb2.MakeRunDeviceTemplateName(deviceId)

	sql := fmt.Sprintf("select sum(duration) as value from %s where status = %d group by ([%s, %s), %s)  align by device",
		deviceKey, status, startTime, endTime, interval)

	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("读取设备运行过程失败:", zap.Error(err))
		return nil, err
	}
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		deviceName := statement.GetText(statement.GetColumnName(0))
		duration := statement.GetDouble(statement.GetColumnName(1))
		_, ok := processList.List[deviceName]
		if !ok {
			processList.List[deviceName] = make([]devicemonitormodel.DeviceStatus, 0)
		}
		processList.List[deviceName] = append(processList.List[deviceName], devicemonitormodel.DeviceStatus{
			Value:  duration,
			Status: status,
			Time:   timestamp,
		})
		continue
	}

	return &processList, nil
}

// ExportTimeData 导入时序数据
func metricEndTime(end uint64) stdtime.Time {
	if end == 0 {
		return stdtime.Now()
	}
	return stdtime.UnixMilli(int64(end))
}

func (i *iotDbExport) ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error {
	if len(data) == 0 {
		return nil
	}

	var samples []iotdb.MetricSample
	for table, list := range data {
		for _, value := range list {
			if value == nil {
				continue
			}
			timestamp := time.MicroToGTime(value.TimeUnixNano).UnixMilli()
			for _, metric := range value.Metrics {
				if metric == nil {
					continue
				}

				var metricValue float64
				if metric.GetValue() == nil {
					metricValue = 0
				} else if _, ok := metric.GetValue().(*v1.TimeDataMetric_AsDouble); ok {
					metricValue = metric.GetAsDouble()
				} else {
					metricValue = float64(metric.GetAsInt())
				}

				samples = append(samples, iotdb.NewMetricSample(table+"."+metric.Field, nil, timestamp, metricValue))
			}
		}
	}
	if len(samples) == 0 {
		return nil
	}

	appender := i.iotDb.Appender()
	if err := appender.Append(samples); err != nil {
		zap.L().Error("iotdb appender append time data error", zap.Error(err))
		return err
	}
	if err := appender.Commit(); err != nil {
		zap.L().Error("iotdb appender commit time data error", zap.Error(err))
		return err
	}
	return nil
}
