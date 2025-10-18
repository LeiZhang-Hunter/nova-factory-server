package systemServiceImpl

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/business/system/systemService"
)

type IDeviceElectricServiceImpl struct {
	dao       systemDao.IDeviceElectricDao
	deviceDao deviceDao.IDeviceDao
}

func NewIDeviceElectricServiceImpl(dao systemDao.IDeviceElectricDao, deviceDao deviceDao.IDeviceDao) systemService.IDeviceElectricService {
	return &IDeviceElectricServiceImpl{
		dao:       dao,
		deviceDao: deviceDao,
	}
}

func (i *IDeviceElectricServiceImpl) Set(c *gin.Context, setting *systemModels.SysDeviceElectricSettingVO) (*systemModels.SysDeviceElectricSetting, error) {
	return i.dao.Set(c, setting)
}
func (i *IDeviceElectricServiceImpl) List(c *gin.Context, req *systemModels.SysDeviceElectricSettingDQL) (*systemModels.SysDeviceElectricSettingData, error) {
	list, err := i.dao.List(c, req)
	if err != nil {
		return list, err
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
	devMap := make(map[uint64]*deviceModels.DeviceVO, 0)
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
		list.Rows[k].DevName = *deviceInfo.Name
	}

	return list, nil
}
func (i *IDeviceElectricServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}

func (i *IDeviceElectricServiceImpl) GetByDeviceId(c *gin.Context, id int64) (*systemModels.SysDeviceElectricSetting, error) {
	return i.dao.GetByDeviceId(c, id)
}

func (i *IDeviceElectricServiceImpl) GetByNoDeviceId(c *gin.Context, id int64, deviceId int64) (*systemModels.SysDeviceElectricSetting, error) {
	return i.dao.GetByNoDeviceId(c, id, deviceId)
}
