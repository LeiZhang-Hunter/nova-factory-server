package metricdaoimpl

import (
	"context"
	"errors"
	"fmt"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/math"
	"sort"
	"strconv"
	"strings"
	stdtime "time"

	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	promlabels "github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"go.uber.org/zap"
)

// prometheusDeviceMetricName 是设备指标写入 Prometheus TSDB 时使用的统一指标名。
// 设备、模板和测点维度保存在 labels 中。
const prometheusDeviceMetricName = "nova_metrics_device"

// iPrometheusTSDBExport 使用本地 Prometheus TSDB 实现 iDaoExport。
type iPrometheusTSDBExport struct {
	txdb iotdb.TSDBStorage
}

type prometheusTSDBDeleter interface {
	Delete(meta iotdb.MetricMeta) error
}

func newPrometheusTSDBExport(tsdb iotdb.TSDBStorage) iDaoExport {
	return &iPrometheusTSDBExport{txdb: tsdb}
}

// prometheusTSDBMetricResult 只收集样本点，用于单指标查询结果。
type prometheusTSDBMetricResult struct {
	points []metricmodels.MetricQueryValue
}

func (r *prometheusTSDBMetricResult) GetKind() string {
	return ""
}

func (r *prometheusTSDBMetricResult) GetProperties() map[string]string {
	return nil
}

// AddSeries 将 Prometheus storage.Series 中的 float 样本转成业务查询点。
func (r *prometheusTSDBMetricResult) AddSeries(series iotdb.MetricSeries) error {
	promSeries, ok := series.(storage.Series)
	if !ok {
		return fmt.Errorf("unsupported prometheus tsdb series %T", series)
	}

	iterator := promSeries.Iterator(nil)
	for valueType := iterator.Next(); valueType != chunkenc.ValNone; valueType = iterator.Next() {
		if valueType != chunkenc.ValFloat {
			continue
		}
		timestamp, value := iterator.At()
		r.points = append(r.points, metricmodels.MetricQueryValue{
			Time:  timestamp,
			Value: value,
		})
	}
	if err := iterator.Err(); err != nil {
		return err
	}
	return nil
}

// prometheusTSDBSeries 保留一个 Prometheus 序列的指标名、labels 和样本点。
type prometheusTSDBSeries struct {
	name   string
	labels map[string]string
	points []metricmodels.MetricQueryValue
}

// prometheusTSDBSeriesResult 收集完整序列信息，用于列表、计数和多序列查询。
type prometheusTSDBSeriesResult struct {
	series []prometheusTSDBSeries
}

func (r *prometheusTSDBSeriesResult) GetKind() string {
	return ""
}

func (r *prometheusTSDBSeriesResult) GetProperties() map[string]string {
	return nil
}

// AddSeries 将 Prometheus 序列和 labels 一起保存下来，便于后续还原设备维度。
func (r *prometheusTSDBSeriesResult) AddSeries(series iotdb.MetricSeries) error {
	promSeries, ok := series.(storage.Series)
	if !ok {
		return fmt.Errorf("unsupported prometheus tsdb series %T", series)
	}

	labels := make(map[string]string)
	promLabels := promSeries.Labels()
	promLabels.Range(func(label promlabels.Label) {
		labels[label.Name] = label.Value
	})

	item := prometheusTSDBSeries{
		name:   labels["__name__"],
		labels: labels,
	}
	iterator := promSeries.Iterator(nil)
	for valueType := iterator.Next(); valueType != chunkenc.ValNone; valueType = iterator.Next() {
		if valueType != chunkenc.ValFloat {
			continue
		}
		timestamp, value := iterator.At()
		item.points = append(item.points, metricmodels.MetricQueryValue{
			Time:  timestamp,
			Value: value,
		})
	}
	if err := iterator.Err(); err != nil {
		return err
	}
	r.series = append(r.series, item)
	return nil
}

// Export 写入设备指标。Prometheus 中统一使用 prometheusDeviceMetricName，
// device_id/template_id/data_id 作为 labels。
func (i *iPrometheusTSDBExport) Export(ctx context.Context, data []*metricmodels.NovaMetricsDevice) error {
	if len(data) == 0 {
		return nil
	}

	samples := make([]iotdb.MetricSample, 0, len(data))
	for _, value := range data {
		if value == nil || value.StartTimeUnix == nil {
			continue
		}
		samples = append(samples, iotdb.NewMetricSample(
			prometheusDeviceMetricName,
			prometheusDeviceMetricProperties(value.DeviceId, value.TemplateId, value.DataId),
			value.StartTimeUnix.UnixMilli(),
			value.Value,
		))
	}
	if len(samples) == 0 {
		return nil
	}

	appender := i.txdb.Appender()
	if err := appender.Append(samples); err != nil {
		zap.L().Error("prometheus tsdb appender append error", zap.Error(err))
		return err
	}
	if err := appender.Commit(); err != nil {
		zap.L().Error("prometheus tsdb appender commit error", zap.Error(err))
		return err
	}
	return nil
}

// ExportTimeData 写入通用时序数据。map 的 key 作为指标名，字段名保存到 field label。
func (i *iPrometheusTSDBExport) ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error {
	if len(data) == 0 {
		return nil
	}

	var samples []iotdb.MetricSample
	for table, list := range data {
		for _, item := range list {
			if item == nil {
				continue
			}
			timestamp := stdtime.Unix(0, int64(item.TimeUnixNano)).UnixMilli()
			for _, metric := range item.Metrics {
				if metric == nil {
					continue
				}

				var value float64
				if metric.GetValue() == nil {
					value = 0
				} else if _, ok := metric.GetValue().(*v1.TimeDataMetric_AsDouble); ok {
					value = metric.GetAsDouble()
				} else {
					value = float64(metric.GetAsInt())
				}

				samples = append(samples, iotdb.NewMetricSample(table, map[string]string{"field": metric.Field}, timestamp, value))
			}
		}
	}
	if len(samples) == 0 {
		return nil
	}

	appender := i.txdb.Appender()
	if err := appender.Append(samples); err != nil {
		zap.L().Error("prometheus tsdb appender append time data error", zap.Error(err))
		return err
	}
	if err := appender.Commit(); err != nil {
		zap.L().Error("prometheus tsdb appender commit time data error", zap.Error(err))
		return err
	}
	return nil
}

// Metric 查询设备单测点数据，并按 req.Step 分钟粒度做平均聚合。
func (i *iPrometheusTSDBExport) Metric(c *gin.Context, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
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

// InstallDevice 是 IoTDB 模板能力，Prometheus TSDB 不需要安装设备模板。
// 这里仅做参数校验，保持接口幂等。
func (i *iPrometheusTSDBExport) InstallDevice(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData) error {
	if deviceId <= 0 {
		return errors.New("设备ID不能为空")
	}
	if device == nil {
		return errors.New("设备配置不能为空")
	}
	return nil
}

// UnInStallDevice 删除 Prometheus TSDB 中指定设备、模板和测点对应的指标序列。
func (i *iPrometheusTSDBExport) UnInStallDevice(c *gin.Context, deviceId int64, templateId int64, dataId int64) error {
	if deviceId <= 0 || templateId <= 0 || dataId <= 0 {
		return errors.New("设备、模板或测点ID不能为空")
	}
	return i.deletePrometheusMetric(iotdb.NewMetricSample(
		prometheusDeviceMetricName,
		prometheusDeviceMetricProperties(uint64(deviceId), uint64(templateId), uint64(dataId)),
		0,
		0,
	))
}

// InstallRunStatusDevice 是 IoTDB 运行状态模板能力，Prometheus TSDB 不需要建模。
func (i *iPrometheusTSDBExport) InstallRunStatusDevice(c *gin.Context, deviceId int64) error {
	if deviceId <= 0 {
		return errors.New("设备ID不能为空")
	}
	return nil
}

// UnInStallRunStatusDevice 删除指定设备运行状态对应的指标序列。
func (i *iPrometheusTSDBExport) UnInStallRunStatusDevice(c *gin.Context, deviceId int64) error {
	if deviceId <= 0 {
		return errors.New("设备ID不能为空")
	}
	return i.deletePrometheusMetric(iotdb.NewMetricSample(iotdb2.MakeRunDeviceTemplateName(deviceId), nil, 0, 0))
}

// Predict 是 IoTDB SQL 扩展能力，Prometheus TSDB 当前返回空结果。
func (i *iPrometheusTSDBExport) Predict(c *gin.Context, deviceId int64, device *devicemodels.SysModbusDeviceConfigData, req *metricmodels.MetricQueryReq) (*metricmodels.MetricQueryData, error) {
	return metricmodels.NewMetricQueryData(), nil
}

// List 查询通用时序数据列表。req.Dev 可传指标名，也可传 指标名.field。
func (i *iPrometheusTSDBExport) List(c *gin.Context, req *devicemonitormodel.DevDataReq) (*devicemonitormodel.DevDataResp, error) {
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
func (i *iPrometheusTSDBExport) Count(c *gin.Context, req *devicemonitormodel.DevDataReq) (uint64, error) {
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
func (i *iPrometheusTSDBExport) Query(c *gin.Context, req *metricmodels.MetricDataQueryReq) (*metricmodels.MetricQueryData, error) {
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

func (i *iPrometheusTSDBExport) deletePrometheusMetric(meta iotdb.MetricMeta) error {
	deleter, ok := i.txdb.(prometheusTSDBDeleter)
	if !ok {
		return errors.New("prometheus tsdb delete is not supported")
	}
	return deleter.Delete(meta)
}

// queryDevDataSeries 将 DevDataReq 中的 dev 条件转换为 Prometheus 序列查询。
func (i *iPrometheusTSDBExport) queryDevDataSeries(startTime, endTime stdtime.Time, devs []string) ([]prometheusTSDBSeries, error) {
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

// queryMetricDataSeries 将 dashboard 查询条件转换为 Prometheus 序列查询。
func (i *iPrometheusTSDBExport) queryMetricDataSeries(startTime, endTime stdtime.Time, req *metricmodels.MetricDataQueryReq) ([]prometheusTSDBSeries, error) {
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

// queryPrometheusSeries 是调用 TSDBStorage Querier 的统一入口。
func (i *iPrometheusTSDBExport) queryPrometheusSeries(startTime, endTime stdtime.Time, meta iotdb.MetricMeta) ([]prometheusTSDBSeries, error) {
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

// prometheusDevDataRows 将 Prometheus 序列样本还原成设备数据列表行。
func prometheusDevDataRows(series []prometheusTSDBSeries) []devicemonitormodel.DevData {
	rows := make([]devicemonitormodel.DevData, 0)
	for _, item := range series {
		dev := item.name
		if field := item.labels["field"]; field != "" {
			dev = item.name + "." + field
		}
		for _, point := range item.points {
			rows = append(rows, devicemonitormodel.DevData{
				Time:       stdtime.UnixMilli(point.Time),
				Value:      point.Value,
				Dev:        dev,
				DeviceID:   prometheusLabelInt64(item.labels["device_id"]),
				TemplateID: prometheusLabelInt64(item.labels["template_id"]),
				DataID:     prometheusLabelInt64(item.labels["data_id"]),
			})
		}
	}
	return rows
}

// prometheusDevDataMetas 生成 List/Count 查询使用的 MetricMeta。
func prometheusDevDataMetas(devs []string) []iotdb.MetricMeta {
	if len(devs) == 0 {
		return []iotdb.MetricMeta{iotdb.NewMetricSample(prometheusDeviceMetricName, nil, 0, 0)}
	}

	metas := make([]iotdb.MetricMeta, 0, len(devs))
	for _, dev := range devs {
		name, field := prometheusSplitMetricField(dev)
		properties := map[string]string{}
		if field != "" {
			properties["field"] = field
		}
		metas = append(metas, iotdb.NewMetricSample(name, properties, 0, 0))
	}
	return metas
}

// prometheusMetricDataMetas 生成 dashboard 查询使用的 MetricMeta。
func prometheusMetricDataMetas(req *metricmodels.MetricDataQueryReq) []iotdb.MetricMeta {
	if len(req.QueryMetric) == 0 {
		name := strings.TrimSpace(req.Name)
		if name == "" {
			name = prometheusDeviceMetricName
		}
		return []iotdb.MetricMeta{iotdb.NewMetricSample(name, nil, 0, 0)}
	}

	metas := make([]iotdb.MetricMeta, 0, len(req.QueryMetric))
	for _, condition := range req.QueryMetric {
		metas = append(metas, iotdb.NewMetricSample(prometheusDeviceMetricName, map[string]string{
			"device_id":   strconv.FormatInt(condition.DeviceId, 10),
			"template_id": strconv.FormatInt(condition.TemplateId, 10),
			"data_id":     strconv.FormatInt(condition.DataId, 10),
		}, 0, 0))
	}
	return metas
}

// prometheusDevDataTimeRange 校验并转换设备数据查询时间范围。
func prometheusDevDataTimeRange(start, end uint64) (stdtime.Time, stdtime.Time, error) {
	if start == 0 {
		return stdtime.Time{}, stdtime.Time{}, errors.New("开始时间不能为空")
	}
	return stdtime.UnixMilli(int64(start)), prometheusMetricEndTime(end), nil
}

// prometheusNormalizeDevDataPage 与 IoTDB 实现保持一致，限制默认分页大小。
func prometheusNormalizeDevDataPage(req *devicemonitormodel.DevDataReq) {
	if req.Size <= 0 {
		req.Size = 20
	}
	if req.Size > 50 {
		req.Size = 50
	}
	if req.Page < 1 {
		req.Page = 1
	}
}

// prometheusMetricInterval 计算 dashboard 聚合间隔，单位为分钟。
func prometheusMetricInterval(start, end uint64, interval int) int {
	if interval != 0 {
		return interval
	}
	calculated := int((end - start) / 60 / 30 / 1000)
	if calculated == 0 {
		return 1
	}
	return calculated
}

// prometheusSplitMetricField 将 指标名.field 拆成指标名和 field label。
func prometheusSplitMetricField(name string) (string, string) {
	name = strings.TrimSpace(name)
	index := strings.LastIndex(name, ".")
	if index <= 0 || index == len(name)-1 {
		return name, ""
	}
	return name[:index], name[index+1:]
}

// prometheusLabelInt64 将 label 值转换为业务 ID，转换失败时返回 0。
func prometheusLabelInt64(value string) int64 {
	parsed, _ := strconv.ParseInt(value, 10, 64)
	return parsed
}

// prometheusDeviceMetricProperties 生成设备指标的 Prometheus labels。
func prometheusDeviceMetricProperties(deviceID, templateID, dataID uint64) map[string]string {
	return map[string]string{
		"device_id":   strconv.FormatUint(deviceID, 10),
		"template_id": strconv.FormatUint(templateID, 10),
		"data_id":     strconv.FormatUint(dataID, 10),
	}
}

// prometheusMetricEndTime 将结束时间转换为 time.Time，未传时使用当前时间。
func prometheusMetricEndTime(end uint64) stdtime.Time {
	if end == 0 {
		return stdtime.Now()
	}
	return stdtime.UnixMilli(int64(end))
}

// bucketAverageMetricValues 按固定时间桶对样本做平均值聚合。
func bucketAverageMetricValues(points []metricmodels.MetricQueryValue, startTime int64, stepMillis int64) []metricmodels.MetricQueryValue {
	if len(points) == 0 {
		return []metricmodels.MetricQueryValue{}
	}
	if stepMillis <= 0 {
		stepMillis = int64(stdtime.Minute / stdtime.Millisecond)
	}

	type bucketValue struct {
		time  int64
		sum   float64
		count int
	}
	buckets := make(map[int64]*bucketValue)
	for _, point := range points {
		bucketIndex := (point.Time - startTime) / stepMillis
		bucketTime := startTime + bucketIndex*stepMillis
		bucket, ok := buckets[bucketTime]
		if !ok {
			bucket = &bucketValue{time: bucketTime}
			buckets[bucketTime] = bucket
		}
		bucket.sum += point.Value
		bucket.count++
	}

	bucketTimes := make([]int64, 0, len(buckets))
	for bucketTime := range buckets {
		bucketTimes = append(bucketTimes, bucketTime)
	}
	sort.Slice(bucketTimes, func(i, j int) bool {
		return bucketTimes[i] < bucketTimes[j]
	})

	values := make([]metricmodels.MetricQueryValue, 0, len(bucketTimes))
	for _, bucketTime := range bucketTimes {
		bucket := buckets[bucketTime]
		values = append(values, metricmodels.MetricQueryValue{
			Time:  bucket.time,
			Value: math.RoundFloat(bucket.sum/float64(bucket.count), 2),
		})
	}
	return values
}
