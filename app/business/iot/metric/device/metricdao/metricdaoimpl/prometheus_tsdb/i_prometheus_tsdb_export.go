package prometheus_tsdb

import (
	"errors"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/math"
	"sort"
	"strconv"
	"strings"
	stdtime "time"

	"github.com/gin-gonic/gin"
)

// prometheusDeviceMetricName 是设备指标写入 Prometheus TSDB 时使用的统一指标名。
// 设备、模板和测点维度保存在 labels 中。
const prometheusDeviceMetricName = "nova_metrics_device"

// iPrometheusTSDBExport 使用本地 Prometheus TSDB 实现 metricdao.IMetricStorageDao。
type iPrometheusTSDBExport struct {
	txdb  iotdb.TSDBStorage
	adder *adder
	query *query
}

func (i *iPrometheusTSDBExport) Questioner() metricdao.IMetricQueryDao {
	//TODO implement me
	return i.query
}

type prometheusTSDBDeleter interface {
	Delete(meta iotdb.MetricMeta) error
}

func NewPrometheusTSDBExport(tsdb iotdb.TSDBStorage) metricdao.IMetricStorageDao {
	return &iPrometheusTSDBExport{
		txdb:  tsdb,
		adder: newAdder(tsdb),
	}
}

// InstallDevice 是 IoTDB 模板能力，Prometheus TSDB 不需要安装设备模板。
func (i *iPrometheusTSDBExport) InstallDevice(c *gin.Context, deviceId int64, templateId int64) error {
	return nil
}

// UnInStallDevice 是 IoTDB 模板能力，Prometheus TSDB 不需要卸载设备模板。
func (i *iPrometheusTSDBExport) UnInStallDevice(c *gin.Context, deviceId int64, templateId int64) error {
	return nil
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

func (i *iPrometheusTSDBExport) deletePrometheusMetric(meta iotdb.MetricMeta) error {
	deleter, ok := i.txdb.(prometheusTSDBDeleter)
	if !ok {
		return errors.New("prometheus tsdb delete is not supported")
	}
	return deleter.Delete(meta)
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

// Template returns the template DAO for the active Prometheus TSDB storage backend.
func (i *iPrometheusTSDBExport) Template() metricdao.IIotStorageTemplateDao {
	return nil
}

// Adder returns the metric adder for the active Prometheus TSDB storage backend.
func (i *iPrometheusTSDBExport) Adder() metricdao.IMetricAdderDao {
	return i.adder
}
