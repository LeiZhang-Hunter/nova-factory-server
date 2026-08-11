package prometheus_tsdb

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	"nova-factory-server/app/datasource/iotdb"
	"sort"
	stdtime "time"
)

type query struct {
	txdb iotdb.TSDBStorage
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

func newQuery(txdb iotdb.TSDBStorage) *query {
	return &query{
		txdb: txdb,
	}
}

// Metric 查询设备单测点数据，并按 req.Step 分钟粒度做平均聚合。
func (i *query) Metric(c *gin.Context, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
	if req == nil {
		return nil, nil
	}
	if req.Start == 0 {
		return nil, errors.New("开始时间不能为空")
	}

	startTime := stdtime.UnixMilli(int64(req.Start))
	endTime := prometheusMetricEndTime(req.End)
	if req.Step <= 0 {
		req.Step = 1
	}

	querier, err := i.txdb.Querier(startTime, endTime)
	if err != nil {
		return nil, err
	}

	result := &prometheusTSDBMetricResult{}
	meta := iotdb.NewMetricSample(
		prometheusDeviceMetricName,
		prometheusDeviceMetricProperties(req.DeviceId, req.TemplateId, req.DataId),
		0,
		0,
	)
	if err := querier.QueryAndClose(meta, nil, result); err != nil {
		zap.L().Error("prometheus tsdb query error", zap.Error(err))
		return nil, err
	}

	data := metricmodels.NewMetricQueryData()
	data.Id = prometheusDeviceMetricName
	data.Values = bucketAverageMetricValues(result.points, startTime.UnixMilli(), int64(req.Step)*int64(stdtime.Minute/stdtime.Millisecond))
	return data, nil
}

// Predict 是 IoTDB SQL 扩展能力，Prometheus TSDB 当前返回空结果。
func (i *query) Predict(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
	return metricmodels.NewMetricQueryData(), nil
}

// List 查询通用时序数据列表。req.Dev 可传指标名，也可传 指标名.field。
func (i *query) List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error) {
	if req == nil {
		return nil, nil
	}
	startTime, endTime, err := prometheusDevDataTimeRange(req.Start, req.End)
	if err != nil {
		return nil, err
	}
	prometheusNormalizeDevDataPage(req)

	series, err := i.queryDevDataSeries(startTime, endTime, req.Dev)
	if err != nil {
		return nil, err
	}
	rows := prometheusDevDataRows(series)
	sort.Slice(rows, func(a, b int) bool {
		return rows[a].Time.After(rows[b].Time)
	})

	start := int((req.Page - 1) * req.Size)
	if start > len(rows) {
		start = len(rows)
	}
	end := start + int(req.Size)
	if end > len(rows) {
		end = len(rows)
	}

	return &devicemonitormodel.DevDataResp{
		Rows:  rows[start:end],
		Total: uint64(len(rows)),
	}, nil
}

// Count 复用 List 的查询条件，统计匹配样本点总数。
func (i *query) Count(c *gin.Context, req *devicemonitormodel.DevDataReq) (uint64, error) {
	if req == nil {
		return 0, nil
	}
	startTime, endTime, err := prometheusDevDataTimeRange(req.Start, req.End)
	if err != nil {
		return 0, err
	}

	series, err := i.queryDevDataSeries(startTime, endTime, req.Dev)
	if err != nil {
		return 0, err
	}
	var total uint64
	for _, item := range series {
		total += uint64(len(item.points))
	}
	return total, nil
}

// Query 用于 dashboard 查询。传 QueryMetric 时按设备维度查询，否则按 req.Name 查询指标名。
func (i *query) Query(c *gin.Context, req *metricmodels.MetricDataQueryReq) (*metricmodels.MetricQueryData, error) {
	if req == nil {
		return nil, nil
	}
	if req.Start == 0 {
		return nil, errors.New("开始时间不能为空")
	}

	startTime := stdtime.UnixMilli(int64(req.Start))
	endTime := prometheusMetricEndTime(req.End)
	interval := prometheusMetricInterval(req.Start, req.End, req.Interval)
	stepMillis := int64(interval) * int64(stdtime.Minute/stdtime.Millisecond)

	series, err := i.queryMetricDataSeries(startTime, endTime, req)
	if err != nil {
		return nil, err
	}

	data := metricmodels.NewMetricQueryData()
	if len(series) <= 1 {
		if len(series) == 1 {
			data.Values = bucketAverageMetricValues(series[0].points, startTime.UnixMilli(), stepMillis)
		}
		return data, nil
	}

	data.MultiValues = make([][]metricmodels.MetricQueryValue, len(series))
	for index, item := range series {
		data.MultiValues[index] = bucketAverageMetricValues(item.points, startTime.UnixMilli(), stepMillis)
	}
	return data, nil
}

// queryDevDataSeries 将 DevDataReq 中的 dev 条件转换为 Prometheus 序列查询。
func (i *query) queryDevDataSeries(startTime, endTime stdtime.Time, devs []string) ([]prometheusTSDBSeries, error) {
	metas := prometheusDevDataMetas(devs)
	series := make([]prometheusTSDBSeries, 0)
	for _, meta := range metas {
		result, err := i.queryPrometheusSeries(startTime, endTime, meta)
		if err != nil {
			return nil, err
		}
		series = append(series, result...)
	}
	return series, nil
}

// queryPrometheusSeries 是调用 TSDBStorage Querier 的统一入口。
func (i *query) queryPrometheusSeries(startTime, endTime stdtime.Time, meta iotdb.MetricMeta) ([]prometheusTSDBSeries, error) {
	querier, err := i.txdb.Querier(startTime, endTime)
	if err != nil {
		return nil, err
	}
	result := &prometheusTSDBSeriesResult{}
	if err := querier.QueryAndClose(meta, nil, result); err != nil {
		return nil, err
	}
	return result.series, nil
}

// queryMetricDataSeries 将 dashboard 查询条件转换为 Prometheus 序列查询。
func (i *query) queryMetricDataSeries(startTime, endTime stdtime.Time, req *metricmodels.MetricDataQueryReq) ([]prometheusTSDBSeries, error) {
	metas := prometheusMetricDataMetas(req)
	series := make([]prometheusTSDBSeries, 0)
	for _, meta := range metas {
		result, err := i.queryPrometheusSeries(startTime, endTime, meta)
		if err != nil {
			return nil, err
		}
		series = append(series, result...)
	}
	return series, nil
}
