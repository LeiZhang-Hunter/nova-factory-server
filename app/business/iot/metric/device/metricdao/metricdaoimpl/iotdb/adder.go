package iotdb

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"go.uber.org/zap"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/time"
)

type adder struct {
	iotDb *iotdb.IotDb
}

func newAdder(iotDb *iotdb.IotDb) *adder {
	return &adder{
		iotDb: iotDb,
	}
}

// Export 导入数据
func (i *adder) Export(ctx context.Context, data []*metricmodels.NovaMetricsDevice) error {
	if len(data) == 0 {
		return nil
	}

	samples := make([]iotdb.MetricSample, 0, len(data))
	for _, value := range data {
		if value == nil || value.StartTimeUnix == nil {
			continue
		}
		name := iotdb2.MakeDeviceDataPath(int64(value.DeviceId), int64(value.DataId)) + ".value"
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

func (i *adder) ExportTimeData(ctx context.Context, data map[string][]*v1.ResourceTimeMetrics) error {
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
