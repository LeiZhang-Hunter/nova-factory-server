package alertServiceImpl

import (
	"errors"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/business/asset/building/buildingDao"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/daemonize/daemonizeDao"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AlertLogServiceImpl struct {
	dao         alertDao.AlertLogClickhouseDao
	ruleDao     alertDao.AlertRuleDao
	agentDao    daemonizeDao.IotAgentDao
	deviceDao   deviceDao.IDeviceDao
	buildingDao buildingDao.BuildingDao
	pointDao    deviceDao.ISysModbusDeviceConfigDataDao
}

func NewAlertLogServiceImpl(dao alertDao.AlertLogClickhouseDao, ruleDao alertDao.AlertRuleDao, agentDao daemonizeDao.IotAgentDao, deviceDao deviceDao.IDeviceDao, buildingDao buildingDao.BuildingDao, pointDao deviceDao.ISysModbusDeviceConfigDataDao) alertService.AlertLogService {
	return &AlertLogServiceImpl{
		dao:         dao,
		ruleDao:     ruleDao,
		agentDao:    agentDao,
		deviceDao:   deviceDao,
		buildingDao: buildingDao,
		pointDao:    pointDao,
	}
}

// Export 导出告警数据
func (log *AlertLogServiceImpl) Export(c *gin.Context, data alertModels.AlertLogData) error {
	// 校验网关帐号密码
	gatewayInfo, err := log.agentDao.GetByObjectId(c, uint64(data.GatewayId))
	if err != nil {
		zap.L().Error("export data error", zap.Error(err))
		return err
	}

	if gatewayInfo == nil {
		return errors.New("gateway not found")
	}

	if gatewayInfo.Username != data.Username || gatewayInfo.Password != data.Password {
		return errors.New("username or password error")
	}

	alertLogList, _ := alertModels.FromDataToNovaAlertLog(&data)
	err = log.dao.Export(c, alertLogList)
	if err != nil {
		return err
	}

	//for _, info := range infos {
	//	alertController.GetAlertRunner().Push(info)
	//}
	return nil
}

func (log *AlertLogServiceImpl) List(c *gin.Context, req *alertModels.SysAlertLogListReq) (*alertModels.NovaAlertLogList, error) {
	list, err := log.dao.List(c, req)
	if err != nil {
		return &alertModels.NovaAlertLogList{
			Rows: make([]*alertModels.NovaAlertLog, 0),
		}, err
	}

	if list == nil {
		return &alertModels.NovaAlertLogList{
			Rows: make([]*alertModels.NovaAlertLog, 0),
		}, nil
	}

	// 收集 ruleIds
	var ruleIdMap map[uint64]bool = make(map[uint64]bool)
	for _, v := range list.Rows {
		ruleIdMap[v.AlertID] = true
	}
	var ruleIds []uint64 = make([]uint64, 0)
	for k := range ruleIdMap {
		ruleIds = append(ruleIds, k)
	}

	rules, err := log.ruleDao.GetByIds(c, ruleIds)
	if err != nil {
		return &alertModels.NovaAlertLogList{
			Rows: make([]*alertModels.NovaAlertLog, 0),
		}, err
	}

	var rulesMap map[uint64]*alertModels.SysAlert = make(map[uint64]*alertModels.SysAlert)
	for _, v := range rules {
		rulesMap[uint64(v.ID)] = v
	}

	// 收集 deviceIds
	var deviceIdSet map[uint64]bool = make(map[uint64]bool)
	for _, v := range list.Rows {
		if v.DeviceId != 0 {
			deviceIdSet[v.DeviceId] = true
		}
	}
	var deviceIds []int64 = make([]int64, 0)
	for k := range deviceIdSet {
		deviceIds = append(deviceIds, int64(k))
	}

	// 批量查询设备
	type deviceInfo struct {
		name         string
		buildingId   uint64
		buildingName string
	}
	deviceMap := make(map[uint64]*deviceInfo)
	if len(deviceIds) > 0 {
		devices, err := log.deviceDao.GetByIds(c, deviceIds)
		if err != nil {
			zap.L().Error("get devices error", zap.Error(err))
		} else {
			for _, d := range devices {
				info := &deviceInfo{}
				if d.Name != nil {
					info.name = *d.Name
				}
				info.buildingId = d.DeviceBuildingId
				deviceMap[d.DeviceId] = info
			}
		}

		// 收集建筑物 ids
		var buildingIdSet map[uint64]bool = make(map[uint64]bool)
		for _, info := range deviceMap {
			if info.buildingId != 0 {
				buildingIdSet[info.buildingId] = true
			}
		}
		var buildingIds []uint64 = make([]uint64, 0)
		for k := range buildingIdSet {
			buildingIds = append(buildingIds, k)
		}

		if len(buildingIds) > 0 {
			buildings, err := log.buildingDao.GetByIds(c, buildingIds)
			if err != nil {
				zap.L().Error("get buildings error", zap.Error(err))
			} else {
				buildingNameMap := make(map[uint64]string)
				for _, b := range buildings {
					buildingNameMap[uint64(b.ID)] = b.Name
				}
				for k, info := range deviceMap {
					if name, ok := buildingNameMap[info.buildingId]; ok {
						deviceMap[k].buildingName = name
					}
				}
			}
		}
	}

	// 收集 DeviceDataID (测点id)
	var pointIdSet map[uint64]bool = make(map[uint64]bool)
	for _, v := range list.Rows {
		if v.DeviceDataID != 0 {
			pointIdSet[v.DeviceDataID] = true
		}
	}
	var pointIds []uint64 = make([]uint64, 0)
	for k := range pointIdSet {
		pointIds = append(pointIds, k)
	}

	pointNameMap := make(map[uint64]string)
	if len(pointIds) > 0 {
		points, err := log.pointDao.GetByIds(c, pointIds)
		if err != nil {
			zap.L().Error("get points error", zap.Error(err))
		} else {
			for _, p := range points {
				pointNameMap[uint64(p.DeviceConfigID)] = p.Name
			}
		}
	}

	// 填充数据
	for k, v := range list.Rows {
		rule, ok := rulesMap[v.AlertID]
		if !ok {
			list.Rows[k].RuleName = "系统告警"
		} else {
			list.Rows[k].RuleName = rule.Name
		}

		if info, ok := deviceMap[v.DeviceId]; ok {
			list.Rows[k].DeviceName = info.name
			list.Rows[k].BuildingName = info.buildingName
		}

		if name, ok := pointNameMap[v.DeviceDataID]; ok {
			list.Rows[k].PointName = name
		}
	}
	return list, nil
}

func (log *AlertLogServiceImpl) Info(c *gin.Context, objectId uint64) (*alertModels.NovaAlertLog, error) {
	info, err := log.dao.GetByObjectId(c, objectId)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	rule, err := log.ruleDao.GetById(c, info.AlertID)
	if err != nil {
		info.RuleName = "系统告警"
		return info, nil
	}
	if rule == nil {
		info.RuleName = "系统告警"
		return info, nil
	}
	info.RuleName = rule.Name
	// 增加设备名称
	deviceInfo, err := log.deviceDao.GetById(c, int64(info.DeviceId))
	if err == nil && deviceInfo != nil {
		info.DeviceName = *deviceInfo.Name
		// 填充建筑物名称
		if deviceInfo.DeviceBuildingId != 0 {
			buildings, bErr := log.buildingDao.GetByIds(c, []uint64{deviceInfo.DeviceBuildingId})
			if bErr == nil && len(buildings) > 0 {
				info.BuildingName = buildings[0].Name
			}
		}
	}

	if err != nil {
		zap.L().Error("get device info error", zap.Error(err))
	}

	// 填充测点名称
	if info.DeviceDataID != 0 {
		point, pErr := log.pointDao.GetById(c, info.DeviceDataID)
		if pErr == nil && point != nil {
			info.PointName = point.Name
		}
	}

	return info, nil
}

// Count 统计clickhouse数据库
func (log *AlertLogServiceImpl) Count(c *gin.Context) (int64, error) {
	return log.dao.Count(c)
}
