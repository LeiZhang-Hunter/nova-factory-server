package prometheus_tsdb

import (
	"fmt"
	promlabels "github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	"nova-factory-server/app/datasource/iotdb"
)

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

func (r *prometheusTSDBSeriesResult) GetName() string {
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
