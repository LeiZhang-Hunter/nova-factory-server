package deviceMonitorServiceImpl

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/datasource/cache"
)

type DeviceMonitorServiceImpl struct {
	dao       deviceDao.IDeviceDao
	metricDao metricDao.IMetricDao
	cache     cache.Cache
}

func NewDeviceMonitorServiceImpl(dao deviceDao.IDeviceDao, cache cache.Cache, metricDao metricDao.IMetricDao) deviceMonitorService.DeviceMonitorService {
	return &DeviceMonitorServiceImpl{
		dao:       dao,
		cache:     cache,
		metricDao: metricDao,
	}
}

func (d *DeviceMonitorServiceImpl) List(c *gin.Context) (*deviceModels.DeviceInfoListData, error) {
	list, err := d.dao.SelectDeviceList(c, &deviceModels.DeviceListReq{})
	if err != nil {
		zap.L().Error("device list select error", zap.Error(err))
		return nil, err
	}

	var deviceIds []uint64 = make([]uint64, 0)
	for _, v := range list.Rows {
		deviceIds = append(deviceIds, v.DeviceId)
	}

	var keys []string = make([]string, 0)
	for _, v := range deviceIds {
		keys = append(keys, device.MakeDeviceKey(uint64(v)))
	}

	slice := d.cache.MGet(c, keys)
	for k, v := range slice.Val() {
		str, ok := v.(string)
		if !ok {
			continue
		}
		if str == "" {
			continue
		}
		deviceMetrics := make(map[uint64]map[uint64]*metricModels.DeviceMetricData) // template_id => data_id
		err := json.Unmarshal([]byte(str), &deviceMetrics)
		if err != nil {
			zap.L().Error("json Unmarshal error", zap.Error(err))
			continue
		}
		list.Rows[k].TemplateList = deviceMetrics
	}
	return list, nil
}

func (d *DeviceMonitorServiceImpl) Metric(c *gin.Context, req *metricModels.MetricQueryReq) (*metricModels.MetricQueryData, error) {
	return d.metricDao.Metric(c, req)
}
