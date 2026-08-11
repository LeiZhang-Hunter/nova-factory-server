package prometheus_tsdb

import (
	"fmt"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	"nova-factory-server/app/datasource/iotdb"
)

// prometheusTSDBMetricResult 只收集样本点，用于单指标查询结果。
type prometheusTSDBMetricResult struct {
	points []metricmodels.MetricQueryValue
}

func (r *prometheusTSDBMetricResult) GetName() string {
	//TODO implement me
	return ""
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
