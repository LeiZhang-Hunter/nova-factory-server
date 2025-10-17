package metricServiceImpl

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricService"
)

type IDevMapServiceImpl struct {
	mapDao    deviceMonitorDao.IDeviceDataReportDao
	deviceDao deviceDao.IDeviceDao
}

func NewIDevMapServiceImpl(mapDao deviceMonitorDao.IDeviceDataReportDao, deviceDao deviceDao.IDeviceDao) metricService.IDevMapService {
	return &IDevMapServiceImpl{
		mapDao:    mapDao,
		deviceDao: deviceDao,
	}
}

// GetDevList iotdb测点名字和设备的映射
func (i *IDevMapServiceImpl) GetDevList(c *gin.Context, dev []string) ([]deviceMonitorModel.SysIotDbDevMapData, error) {
	mapList, err := i.mapDao.GetDevList(c, dev)
	if err != nil {
		return mapList, err
	}
	if len(mapList) == 0 {
		return mapList, nil
	}
	deviceIdMap := make(map[int64]int)
	deviceIds := make([]int64, 0)
	for _, v := range mapList {
		deviceIdMap[v.DeviceID] = 0
	}

	for id, _ := range deviceIdMap {
		if id != 0 {
			deviceIds = append(deviceIds, id)
		}
	}

	deviceList, err := i.deviceDao.GetByIds(c, deviceIds)
	if err != nil {
		zap.L().Error("get dev list error", zap.Error(err))
	}

	if len(deviceList) == 0 {
		return mapList, nil
	}
	devMap := make(map[uint64]*deviceModels.DeviceVO, 0)
	for _, v := range deviceList {
		devMap[(v.DeviceId)] = v
	}

	for k, v := range mapList {
		deviceInfo, ok := devMap[uint64(v.DeviceID)]
		if !ok {
			continue
		}
		if deviceInfo.Name == nil {
			continue
		}
		mapList[k].DevName = *deviceInfo.Name
	}

	return mapList, nil
}
