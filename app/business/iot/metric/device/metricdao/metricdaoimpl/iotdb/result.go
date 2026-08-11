package iotdb

import (
	"fmt"
	metricmodels "nova-factory-server/app/business/iot/metric/device/metricmodels/entity"
	"nova-factory-server/app/datasource/iotdb"
	"nova-factory-server/app/utils/math"
)

type iotMetricQueryResult struct {
	data *metricmodels.MetricQueryData
}

func (r *iotMetricQueryResult) GetName() string {
	//TODO implement me
	return ""
}

func (r *iotMetricQueryResult) GetKind() string {
	return ""
}

func (r *iotMetricQueryResult) GetProperties() map[string]string {
	return nil
}

func (r *iotMetricQueryResult) AddSeries(series iotdb.MetricSeries) error {
	iotSeries, ok := series.(iotdb.IotDBMetricSeries)
	if !ok {
		return fmt.Errorf("unsupported iotdb metric series %T", series)
	}
	for _, sample := range iotSeries.Samples {
		r.data.Values = append(r.data.Values, metricmodels.MetricQueryValue{
			Time:  sample.Timestamp,
			Value: math.RoundFloat(sample.Value, 2),
		})
	}
	return nil
}
