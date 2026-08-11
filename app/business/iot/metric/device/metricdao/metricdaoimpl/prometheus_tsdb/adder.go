package prometheus_tsdb

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"go.uber.org/zap"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	"nova-factory-server/app/datasource/iotdb"
	stdtime "time"
)

type adder struct {
	txdb iotdb.TSDBStorage
}

func newAdder(txdb iotdb.TSDBStorage) *adder {
	return &adder{
		txdb: txdb,
	}
}

// Export 写入设备指标。Prometheus 中统一使用 prometheusDeviceMetricName，
// device_id/template_id/data_id 作为 labels。
func (i *adder) Export(ctx context.Context, data []*metricmodels.NovaMetricsDevice) error {
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
func (i *adder) ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error {
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
