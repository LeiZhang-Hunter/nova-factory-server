package metricServiceImpl

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/business/metric/device/metricService"
	time2 "nova-factory-server/app/utils/time"
)

type IControlLogServiceImpl struct {
	dao metricDao.IControlLogDao
}

func NewIControlLogServiceImpl(dao metricDao.IControlLogDao) metricService.IControlLogService {
	return &IControlLogServiceImpl{
		dao: dao,
	}
}

func (i *IControlLogServiceImpl) Export(ctx context.Context, request *v1.ExportControlLogRequest) error {

	if request == nil {
		return nil
	}

	if len(request.ResourceMetrics) == 0 {
		return nil
	}

	var data []*metricModels.NovaControlLog = make([]*metricModels.NovaControlLog, 0)

	for _, resourceMetric := range request.ResourceMetrics {
		if resourceMetric == nil {
			continue
		}

		if len(resourceMetric.GetMetrics()) == 0 {
			continue
		}

		for _, v := range resourceMetric.GetMetrics() {
			data = append(data, &metricModels.NovaControlLog{
				DeviceId:      v.DeviceId,
				DataId:        v.DataId,
				Message:       v.Message,
				Type:          v.Type,
				SeriesId:      v.GetStartTimeUnixNano(),
				Attributes:    make(map[string]string),
				StartTimeUnix: time2.MicroToGTime(v.StartTimeUnixNano),
				TimeUnix:      time2.MicroToGTime(v.TimeUnixNano),
			})
		}
	}

	return i.dao.Export(ctx, data)
}
