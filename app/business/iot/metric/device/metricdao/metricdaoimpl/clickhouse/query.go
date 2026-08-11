package clickhouse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	"nova-factory-server/app/datasource/clickhouse"
	"nova-factory-server/app/utils/math"
	"nova-factory-server/app/utils/time"
)

type query struct {
	tableName  string
	clickhouse *clickhouse.ClickHouse
}

func (i *query) CounterByTimeRange(startTime int64, endTime int64, interval string) (*metricmodels.MetricQueryData, error) {
	//TODO implement me
	return nil, nil
}

func (i *query) CounterByDevice(c *gin.Context, startTime int64, endTime int64, limit int) (*devicemonitormodel.TypeDeviceCounterRank, error) {
	//TODO implement me
	return nil, nil
}

func (i *query) StatDeviceStatus(c *gin.Context, startTime string, endTime string, status int) (*devicemonitormodel.DeviceStatusList, error) {
	//TODO implement me
	return nil, nil
}

func (i *query) StatDeviceProcess(c *gin.Context, startTime string, endTime string, interval string, status int) (*devicemonitormodel.DeviceProcessList, error) {
	//TODO implement me
	return nil, nil
}

func (i *query) StatDeviceRunStatus(c *gin.Context, startTime string, endTime string) ([]devicemonitormodel.DeviceRunStat, error) {
	//TODO implement me
	return nil, nil
}

func (i *query) StatDeviceStatusByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, status int) (*devicemonitormodel.DeviceStatusList, error) {
	//TODO implement me
	return nil, nil
}

func (i *query) StatDeviceProcessByDeviceId(c *gin.Context, startTime string, endTime string, deviceId int64, interval string, status int) (*devicemonitormodel.DeviceProcessList, error) {
	//TODO implement me
	return nil, nil
}

func newQuery(tableName string, clickhouse *clickhouse.ClickHouse) *query {
	return &query{
		tableName:  tableName,
		clickhouse: clickhouse,
	}
}

func (i *query) Metric(c *gin.Context, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
	if req == nil {
		return nil, nil
	}

	model := i.clickhouse.DB().Table(i.tableName)

	var startTime string
	if req.Start > 0 {
		startTime = time.GetStartTime(req.Start, 200)
	}
	endTime := time.GetEndTimeUseNow(req.End, true)

	if startTime != "" && endTime != "" {

		model = model.Where("time_unix >= ?", startTime)
	}
	if endTime != "" {
		model = model.Where("time_unix <= ?", endTime)
	}
	if req.Step <= 0 {
		req.Step = 1
	}

	var list []*metricmodels.NovaMetricsDevice = make([]*metricmodels.NovaMetricsDevice, 0)
	model = model.
		Select(fmt.Sprintf(
			"`toStartOfInterval`(`time_unix`, INTERVAL %d minute) AS `time_unix`, %s as value",
			req.Step, "AVG(value)"))

	ret := model.Where("device_id = ?", req.DeviceId).
		Where("template_id = ?", req.TemplateId).
		Where("data_id = ?", req.DataId).
		Group("time_unix").Order("time_unix asc").Limit(500).Find(&list)
	if ret.Error != nil {
		return nil, ret.Error
	}

	if len(list) == 0 {
		return &metricmodels.MetricQueryData{
			Labels: make(map[string]string),
			Values: make([]metricmodels.MetricQueryValue, len(list)),
		}, nil
	}

	var data metricmodels.MetricQueryData
	data.Labels = make(map[string]string)
	data.Values = make([]metricmodels.MetricQueryValue, len(list))
	for i := 0; i < len(list); i++ {
		data.Values[i].Time = list[i].TimeUnix.UnixMilli()
		data.Values[i].Value = math.RoundFloat(list[i].Value, 2)
	}
	return &data, nil
}

func (i *query) Predict(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
	return &metricmodels.MetricQueryData{
		Labels: make(map[string]string),
		Values: make([]metricmodels.MetricQueryValue, 0),
	}, nil
}

func (i *query) List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error) {
	return nil, nil
}

func (i *query) Count(c *gin.Context, req *devicemonitormodel.DevDataReq) (uint64, error) {
	return 0, nil
}

func (i *query) Query(c *gin.Context, req *metricmodels.MetricDataQueryReq) (*metricmodels.MetricQueryData, error) {
	return nil, nil
}
