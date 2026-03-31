package deviceMonitorServiceImpl

import (
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitordao"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitorservice"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IDeviceDataReportServiceImpl struct {
	dao       devicemonitordao.IDeviceDataReportDao
	deviceDao devicedao.IDeviceDao
}

func NewIDeviceDataReportServiceImpl(dao devicemonitordao.IDeviceDataReportDao, deviceDao devicedao.IDeviceDao) devicemonitorservice.IDeviceDataReportService {
	return &IDeviceDataReportServiceImpl{
		dao:       dao,
		deviceDao: deviceDao,
	}
}

func (i *IDeviceDataReportServiceImpl) DevList(c *gin.Context) ([]devicemonitormodel.SysIotDbDevMapData, error) {
	list, err := i.dao.DevList(c)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return list, nil
	}

	deviceIdMap := make(map[int64]int)
	deviceIds := make([]int64, 0)
	for _, v := range list {
		deviceIdMap[v.DeviceID] = 0
	}

	for id := range deviceIdMap {
		if id != 0 {
			deviceIds = append(deviceIds, id)
		}
	}

	deviceList, err := i.deviceDao.GetByIds(c, deviceIds)
	if err != nil {
		zap.L().Error("get dev list error", zap.Error(err))
	}

	if len(deviceList) == 0 {
		return list, nil
	}
	devMap := make(map[uint64]*devicemodels.DeviceVO, 0)
	for _, v := range deviceList {
		devMap[(v.DeviceId)] = v
	}

	for k, v := range list {
		deviceInfo, ok := devMap[uint64(v.DeviceID)]
		if !ok {
			continue
		}
		if deviceInfo.Name == nil {
			continue
		}
		list[k].DevName = *deviceInfo.Name
	}
	return list, nil
}

func (i *IDeviceDataReportServiceImpl) GetDevList(c *gin.Context, req *devicemonitormodel.DevListReq) (*devicemonitormodel.DevListResp, error) {
	list, err := i.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return list, nil
	}
	if len(list.Rows) == 0 {
		return list, nil
	}

	deviceIdMap := make(map[int64]int)
	deviceIds := make([]int64, 0)
	for _, v := range list.Rows {
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
		return list, nil
	}
	devMap := make(map[uint64]*devicemodels.DeviceVO, 0)
	for _, v := range deviceList {
		devMap[(v.DeviceId)] = v
	}

	for k, v := range list.Rows {
		deviceInfo, ok := devMap[uint64(v.DeviceID)]
		if !ok {
			continue
		}
		if deviceInfo.Name == nil {
			continue
		}
		list.Rows[k].DeviceName = *deviceInfo.Name
	}
	return list, nil
}

func (i *IDeviceDataReportServiceImpl) SetDevMap(c *gin.Context, info *devicemonitormodel.SetDevMapInfo) error {
	return i.dao.Save(c, &devicemonitormodel.SysIotDbDevMap{
		ID:       info.ID,
		Device:   info.Device,
		DataName: info.DataName,
		Unit:     info.Unit,
	})

}

func (i *IDeviceDataReportServiceImpl) RemoveDevMap(c *gin.Context, dev string) error {
	return i.dao.Remove(c, dev)

}
