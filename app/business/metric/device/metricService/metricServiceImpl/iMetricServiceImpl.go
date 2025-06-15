package metricServiceImpl

import (
	"context"
	"encoding/json"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"go.uber.org/zap"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/business/metric/device/metricService"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/datasource/cache"
	time2 "nova-factory-server/app/utils/time"
	"time"
)

type IMetricServiceImpl struct {
	dao   metricDao.IMetricDao
	cache cache.Cache
}

func NewIMetricServiceImpl(dao metricDao.IMetricDao, cache cache.Cache) metricService.IMetricService {
	return &IMetricServiceImpl{
		dao:   dao,
		cache: cache,
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
					DataId:        mVale.DataId,
					SeriesId:      request.SeriesId,
					StartTimeUnix: time2.MicroToGTime(mVale.StartTimeUnixNano),
					TimeUnix:      time2.MicroToGTime(mVale.TimeUnixNano),
					Value:         v,
				})
			}
			err := m.export(c, values, request.SeriesId)
			if err != nil {
				zap.L().Error("export error", zap.Error(err))
				continue
			}
		}

	}
	return nil
}

func (m *IMetricServiceImpl) export(c context.Context, values []*metricModels.NovaMetricsDevice, seriesId uint64) error {
	data := metricModels.NewMetricMap()

	// v.Value
	// 格式化数据结构
	for _, v := range values {
		_, ok := data.Data[uint64(v.DeviceId)]
		if !ok {
			data.Data[uint64(v.DeviceId)] = make(map[uint64]map[uint64]*metricModels.DeviceMetricData)
		}

		_, ok = data.Data[uint64(v.DeviceId)][uint64(v.TemplateId)]
		if !ok {
			data.Data[uint64(v.DeviceId)][v.TemplateId] = make(map[uint64]*metricModels.DeviceMetricData)
		}

		mapValue := make(map[string]string)
		data.Data[uint64(v.DeviceId)][uint64(v.TemplateId)][uint64(v.DataId)] = &metricModels.DeviceMetricData{
			SeriesId:      seriesId,
			Attributes:    mapValue,
			StartTimeUnix: v.StartTimeUnix,
			TimeUnix:      v.TimeUnix,
			Value:         v.Value,
		}

		bytes, err := json.Marshal(data.Data[uint64(v.DeviceId)])
		if err != nil {
			zap.L().Error("json marshal error", zap.Error(err))
			continue
		}
		m.cache.Set(c, device.MakeDeviceKey(uint64(v.DeviceId)), string(bytes), 600*time.Second)
	}

	return m.dao.Export(c, values)
}
