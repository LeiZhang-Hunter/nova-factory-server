package metricDaoIMpl

import (
	"context"
	"errors"
	"fmt"
	"github.com/apache/iotdb-client-go/client"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricModels"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/time"
	"strings"
)

type iotDbExport struct {
	iotDb *iotdb.IotDb
}

func newIotDbExport(iotDb *iotdb.IotDb) iDaoExport {
	return &iotDbExport{
		iotDb: iotDb,
	}
}

func (i *iotDbExport) Export(ctx context.Context, data []*metricModels.NovaMetricsDevice) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		return err
	}
	defer i.iotDb.PutSession(session)

	// 变更数据结构，根据设备id分类
	var dataMap map[uint64][]*metricModels.NovaMetricsDevice = make(map[uint64][]*metricModels.NovaMetricsDevice)
	for _, d := range data {
		_, ok := dataMap[d.DeviceId]
		if !ok {
			dataMap[d.DeviceId] = make([]*metricModels.NovaMetricsDevice, 0)
		}
		dataMap[d.DeviceId] = append(dataMap[d.DeviceId], d)
	}

	type pointData struct {
		measurementsSlice [][]string
		dataTypes         [][]client.TSDataType
		values            [][]interface{}
		timestamps        []int64
	}

	var pointDataMap map[string]*pointData = make(map[string]*pointData)

	for deviceId, list := range dataMap {
		for _, value := range list {
			name := iotdb2.MakeDeviceTemplateName(int64(deviceId), int64(value.TemplateId), int64(value.DataId))
			_, ok := pointDataMap[name]
			if !ok {
				pointDataMap[name] = &pointData{
					measurementsSlice: make([][]string, 0),
					dataTypes:         [][]client.TSDataType{},
					values:            [][]interface{}{},
					timestamps:        []int64{},
				}
			}
			now := value.StartTimeUnix.UnixMilli()
			pointDataMap[name].measurementsSlice = append(pointDataMap[name].measurementsSlice, []string{
				"value",
			})
			pointDataMap[name].dataTypes = append(pointDataMap[name].dataTypes, []client.TSDataType{client.DOUBLE})
			pointDataMap[name].values = append(pointDataMap[name].values, []interface{}{(value.Value)})
			pointDataMap[name].timestamps = append(pointDataMap[name].timestamps, now)
		}
	}

	for name, point := range pointDataMap {
		_, err := session.InsertRecordsOfOneDevice(name, point.timestamps, point.measurementsSlice, point.dataTypes, point.values, false)
		if err != nil {
			zap.L().Error("InsertRecordsOfOneDevice error", zap.Error(err))
			continue
		}
	}

	return nil
}

func (i *iotDbExport) Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error) {
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
	if req.Step <= 0 {
		req.Step = 1
	}
	name := iotdb2.MakeDeviceTemplateName(int64(req.DeviceId), int64(req.TemplateId), int64(req.DataId))
	var timeout int64 = 5000
	var data *metricModels.MetricQueryData = metricModels.NewMetricQueryData()
	sql := fmt.Sprintf("select avg(value) as value from %s group by([%s, %s), %dm, %dm);",
		name, startTime, endTime, req.Step, req.Step)
	fmt.Println(sql)
	// select avg(value) from root.device.dev375986234780028928 group by([2025-07-07 20:52:28, 2025-07-07 21:52:28), 3m, 3m);
	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		v := statement.GetDouble("value")
		data.Values = append(data.Values, metricModels.MetricQueryValue{
			Time:  timestamp,
			Value: fmt.Sprintf("%f", v),
		})

	}
	data.Id = name
	return data, nil

}

// InstallDevice 安装设备模板
func (i *iotDbExport) InstallDevice(c *gin.Context, deviceId int64, device *deviceModels.SysModbusDeviceConfigData) error {
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

// UnInStallDevice 卸载设备模板
func (i *iotDbExport) UnInStallDevice(c *gin.Context, deviceId int64, templateId int64, dataId int64) error {
	session, err := i.iotDb.GetSession()
	if err != nil {
		zap.L().Error("读取session失败", zap.Error(err))
		return err
	}
	defer i.iotDb.PutSession(session)

	name := fmt.Sprintf(iotdb2.ROOT_DEVICE_TEMPLATE_NAME, deviceId, templateId, dataId)

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
func (i *iotDbExport) Predict(c *gin.Context, deviceId int64, device *deviceModels.SysModbusDeviceConfigData, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error) {
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
	var data *metricModels.MetricQueryData = metricModels.NewMetricQueryData()
	if device.AggFunction == "" {
		device.AggFunction = "avg"
	}
	sql := fmt.Sprintf("select %s(value) as value from %s group by([%s, %s), %dm, %dm)", device.AggFunction,
		name, startTime, endTime, req.Step, req.Step)
	predictSql := fmt.Sprintf("call inference(_STLForecaster, \"%s\", generateTime=True, predict_length=10)",
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
		data.Values = append(data.Values, metricModels.MetricQueryValue{
			Time:  timestamp,
			Value: fmt.Sprintf("%f", v),
		})

	}
	data.Id = name
	return data, nil
}

func (i *iotDbExport) List(c *gin.Context, req *deviceMonitorModel.DevDataReq) (*deviceMonitorModel.DevDataResp, error) {
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
	var resp deviceMonitorModel.DevDataResp
	resp.Rows = make([]deviceMonitorModel.DevData, 0)
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
		resp.Rows = append(resp.Rows, deviceMonitorModel.DevData{
			Time:  time.MillToTime(timestamp),
			Value: value,
			Dev:   device,
		})

	}

	return &resp, nil
}

func (i *iotDbExport) Count(c *gin.Context, req *deviceMonitorModel.DevDataReq) (uint64, error) {
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
	var resp deviceMonitorModel.DevDataResp
	resp.Rows = make([]deviceMonitorModel.DevData, 0)
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
func (i *iotDbExport) Query(c *gin.Context, req *metricModels.MetricDataQueryReq) (*metricModels.MetricQueryData, error) {
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
	if req.Step <= 0 {
		req.Step = 1
	}
	var timeout int64 = 5000
	var data *metricModels.MetricQueryData = metricModels.NewMetricQueryData()
	var sql string
	if req.Step != 0 {
		sql = fmt.Sprintf("select %s from %s group by([%s, %s), %dm, %dm);",
			req.Expression, req.Name, startTime, endTime, req.Interval, req.Step)
	} else {
		sql = fmt.Sprintf("select %s from %s group by([%s, %s), %dm);",
			req.Expression, req.Name, startTime, endTime, req.Interval)
	}

	if req.Predict.Enable {
		sql = fmt.Sprintf("call inference(%s, \"%s\", generateTime=True, predict_length=10)",
			req.Predict.Model, req.Predict.Param, sql)
	}

	// select avg(value) from root.device.dev375986234780028928 group by([2025-07-07 20:52:28, 2025-07-07 21:52:28), 3m, 3m);
	statement, err := session.ExecuteQueryStatement(sql, &timeout)
	if err != nil {
		zap.L().Error("ExecuteQueryStatement error", zap.Error(err))
		return nil, err
	}
	for next, err := statement.Next(); err == nil && next; next, err = statement.Next() {
		timestamp := statement.GetTimestamp()
		v := statement.GetDouble("value")
		data.Values = append(data.Values, metricModels.MetricQueryValue{
			Time:  timestamp,
			Value: fmt.Sprintf("%f", v),
		})

	}
	return data, nil
}
