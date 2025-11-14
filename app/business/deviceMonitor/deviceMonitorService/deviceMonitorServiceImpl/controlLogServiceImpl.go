package deviceMonitorServiceImpl

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
)

type ControlLogServiceImpl struct {
	dao           metricDao.IControlLogDao
	deviceDao     deviceDao.IDeviceDao
	deviceDataDao deviceDao.ISysModbusDeviceConfigDataDao
}

func NewControlLogServiceImpl(dao metricDao.IControlLogDao, deviceDao deviceDao.IDeviceDao, deviceDataDao deviceDao.ISysModbusDeviceConfigDataDao) deviceMonitorService.ControlLogService {
	return &ControlLogServiceImpl{
		dao:           dao,
		deviceDao:     deviceDao,
		deviceDataDao: deviceDataDao,
	}
}

func (i *ControlLogServiceImpl) List(c *gin.Context, req *deviceMonitorModel.ControlLogListReq) (*metricModels.NovaControlLogList, error) {
	list, err := i.dao.List(c, req)
	if err != nil {
		zap.L().Error("get control log list error", zap.Error(err))
		return nil, err
	}

	if list == nil {
		return &metricModels.NovaControlLogList{
			Rows:  []*metricModels.NovaControlLog{},
			Total: 0,
		}, nil
	}

	if len(list.Rows) == 0 {
		return &metricModels.NovaControlLogList{
			Rows:  []*metricModels.NovaControlLog{},
			Total: 0,
		}, nil
	}

	// 设备id集合
	deviceIdMap := make(map[uint64]int)
	deviceIds := make([]int64, 0)
	dataIdMap := make(map[uint64]int)
	dataIds := make([]uint64, 0)

	for _, v := range list.Rows {
		deviceIdMap[v.DeviceId] = 0
		dataIdMap[v.DataId] = 0
	}

	for id, _ := range deviceIdMap {
		if id != 0 {
			deviceIds = append(deviceIds, int64(id))
		}
	}

	for id, _ := range dataIdMap {
		if id != 0 {
			dataIds = append(dataIds, uint64(id))
		}
	}

	// 读取设备列表
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

	// 读取测点列表
	datas, err := i.deviceDataDao.GetByIds(c, dataIds)
	if err != nil {
		zap.L().Error("get dev list error", zap.Error(err))
		return nil, err
	}
	devDataMap := make(map[uint64]*deviceModels.SysModbusDeviceConfigData)
	for _, v := range datas {
		devDataMap[uint64(v.DeviceConfigID)] = v
	}

	for k, v := range list.Rows {
		deviceInfo, ok := devMap[(v.DeviceId)]
		if !ok {
			continue
		}
		if deviceInfo == nil {
			continue
		}
		if deviceInfo.Name != nil {
			list.Rows[k].DeviceName = *deviceInfo.Name
		} else {
			list.Rows[k].DeviceName = ""
		}

		deviceDataInfo, ok := devDataMap[(v.DataId)]
		if !ok {
			continue
		}
		if deviceDataInfo == nil {
			continue
		}

		list.Rows[k].DataName = deviceDataInfo.Name
	}
	return list, nil
}
