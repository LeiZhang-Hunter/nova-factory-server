package iotdb

import (
	"errors"
	"fmt"
	"github.com/apache/iotdb-client-go/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/math"
	"nova-factory-server/app/utils/time"
	"strings"
	stdtime "time"
)

type query struct {
	iotDb *iotdb.IotDb
}

func newQuery(iotDb *iotdb.IotDb) *query {
	return &query{
		iotDb: iotDb,
	}
}

func (i *query) Metric(c *gin.Context, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
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

	name := iotdb2.MakeDeviceDataPath(int64(req.DeviceId), int64(req.DataId))
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

// Predict 趋势预测
func (i *query) Predict(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
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
	name := iotdb2.MakeDeviceDataPath(int64(req.DeviceId), int64(req.DataId))
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

func (i *query) List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error) {
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

func (i *query) Count(c *gin.Context, req *devicemonitormodel.DevDataReq) (uint64, error) {
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
func (i *query) Query(c *gin.Context, req *metricmodels.MetricDataQueryReq) (*metricmodels.MetricQueryData, error) {
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

func (i *query) CounterByTimeRange(startTime int64, endTime int64, interval string) (*metricmodels.MetricQueryData, error) {
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

func (i *query) CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*devicemonitormodel.TypeDeviceCounterRank, error) {
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

func (i *query) StatDeviceStatus(c *gin.Context, startTime string, endTime string,
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

func (i *query) StatDeviceProcess(c *gin.Context, startTime string, endTime string, interval string,
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

func (i *query) StatDeviceRunStatus(c *gin.Context, startTime string,
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

func (i *query) StatDeviceStatusByDeviceId(c *gin.Context, startTime string, endTime string,
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

func (i *query) StatDeviceProcessByDeviceId(c *gin.Context, startTime string, endTime string,
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
