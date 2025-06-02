package metricServiceImpl

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/business/metric/device/metricService"
	"nova-factory-server/app/utils"
)

type IMetricServiceImpl struct {
	dao metricDao.IMetricDao
}

func NewIMetricServiceImpl(dao metricDao.IMetricDao) metricService.IMetricService {
	return &IMetricServiceImpl{
		dao: dao,
	}
}

func (m *IMetricServiceImpl) Export(c context.Context, request *v1.ExportMetricsServiceRequest) error {
	if request == nil {
		return nil
	}
	if len(request.ResourceMetrics) == 0 {
		return nil
	}

	for _, resourceMetric := range request.ResourceMetrics {
		if resourceMetric == nil {
			continue
		}

		if len(resourceMetric.ScopeMetrics) == 0 {
			continue
		}

		for _, metric := range resourceMetric.ScopeMetrics {
			if metric == nil {
				continue
			}

			if len(metric.Metrics) == 0 {
				continue
			}

			var values []*metricModels.NovaMetricsDevice = make([]*metricModels.NovaMetricsDevice, 0)
			for _, mVale := range metric.Metrics {
				var v float64
				if mVale.GetValue() == nil {
					v = 0
				} else {
					ptr, ok := mVale.GetValue().(*v1.DeviceMetric_AsDouble)
					if ok {
						v = ptr.AsDouble
					} else {
						v = float64(mVale.GetAsInt())
					}
				}
				values = append(values, &metricModels.NovaMetricsDevice{
					DeviceId:   mVale.DeviceId,
					TemplateId: mVale.TemplateId,
					Attributes: map[string]string{
						"test": "1111",
					},
					StartTimeUnix: utils.NanoToGTime(mVale.StartTimeUnixNano),
					TimeUnix:      utils.NanoToGTime(mVale.TimeUnixNano),
					Value:         v,
				})
			}
			err := m.export(c, values)
			if err != nil {
				return err
			}

		}

	}
	return nil
}

func (m *IMetricServiceImpl) export(c context.Context, values []*metricModels.NovaMetricsDevice) error {
	return m.dao.Export(c, values)
}
