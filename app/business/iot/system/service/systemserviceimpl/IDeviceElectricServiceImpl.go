package systemserviceimpl

import (
	"errors"
	"fmt"
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitordao"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	"nova-factory-server/app/business/iot/system/dao"
	"nova-factory-server/app/business/iot/system/models"
	"nova-factory-server/app/business/iot/system/service"
	"nova-factory-server/app/constant/iotdb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IDeviceElectricServiceImpl struct {
	dao       dao.IDeviceElectricDao
	deviceDao devicedao.IDeviceDao
	metricDao metricdao.IMetricDao
	mapDao    devicemonitordao.IDeviceDataReportDao
}

func NewIDeviceElectricServiceImpl(dao dao.IDeviceElectricDao, deviceDao devicedao.IDeviceDao,
	metricDao metricdao.IMetricDao, mapDao devicemonitordao.IDeviceDataReportDao) service.IDeviceElectricService {
	return &IDeviceElectricServiceImpl{
		dao:       dao,
		deviceDao: deviceDao,
		metricDao: metricDao,
		mapDao:    mapDao,
	}
}

func (i *IDeviceElectricServiceImpl) Set(c *gin.Context, setting *models.SysDeviceElectricSettingVO) (*models.SysDeviceElectricSetting, error) {
	if setting == nil {
		return nil, errors.New("setting is nil")
	}
	err := i.metricDao.InstallRunStatusDevice(c, setting.DeviceID)
	if err != nil {
		zap.L().Error("install device run status dev table error", zap.Error(err))
	}

	devKey := iotdb.MakeRunDeviceTemplateName(setting.DeviceID)
	err = i.mapDao.Save(c, &devicemonitormodel.SysIotDbDevMap{
		DeviceID:   setting.DeviceID,
		TemplateID: 0,
		DataID:     0,
		Device:     devKey,
		DataName:   setting.Name,
		Unit:       "",
	})
	if err != nil {
		zap.L().Error("save iotdb device map error", zap.Error(err))
	}
	return i.dao.Set(c, setting)
}
func (i *IDeviceElectricServiceImpl) List(c *gin.Context, req *models.SysDeviceElectricSettingDQL) (*models.SysDeviceElectricSettingData, error) {
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
		list.Rows[k].DevName = *deviceInfo.Name
	}

	return list, nil
}
func (i *IDeviceElectricServiceImpl) Remove(c *gin.Context, ids []string) error {
	list, err := i.dao.GetByIds(c, ids)
	if err != nil {
		zap.L().Error("get dev list error", zap.Error(err))
		return err
	}
	if len(list) != 0 {
		for _, v := range list {
			err = i.metricDao.UnInStallRunStatusDevice(c, v.DeviceID)
			if err != nil {
				zap.L().Error("uninstall device run status dev table error", zap.Error(err))
			}
		}
	}
	fmt.Print(list)
	return i.dao.Remove(c, ids)
}

func (i *IDeviceElectricServiceImpl) GetByDeviceId(c *gin.Context, id int64) (*models.SysDeviceElectricSetting, error) {
	return i.dao.GetByDeviceId(c, id)
}

func (i *IDeviceElectricServiceImpl) GetByNoDeviceId(c *gin.Context, id int64, deviceId int64) (*models.SysDeviceElectricSetting, error) {
	return i.dao.GetByNoDeviceId(c, id, deviceId)
}
