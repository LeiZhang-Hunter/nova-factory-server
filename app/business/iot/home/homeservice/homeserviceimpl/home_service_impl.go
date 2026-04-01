package homeserviceimpl

import (
	"errors"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/admin/monitor/monitormodels"
	"nova-factory-server/app/business/iot/alert/alertmodels"
	"nova-factory-server/app/business/iot/alert/alertservice"
	"nova-factory-server/app/business/iot/asset/building/buildingdao"
	"nova-factory-server/app/business/iot/asset/device/deviceservice"
	"nova-factory-server/app/business/iot/craft/craftrouteservice"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitordao"
	"nova-factory-server/app/business/iot/home/homemodels"
	"nova-factory-server/app/business/iot/home/homeservice"
	"nova-factory-server/app/datasource/cache"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type HomeServiceImpl struct {
	buildingDao          buildingdao.BuildingDao
	deviceService        deviceservice.IDeviceService
	cache                cache.Cache
	craftService         craftrouteservice.ICraftRouteService
	alertRuleService     alertservice.AlertLogService
	deviceMonitorCalcDao devicemonitordao.DeviceMonitorCalcDao
}

func NewHomeServiceImpl(buildingDao buildingdao.BuildingDao, deviceService deviceservice.IDeviceService,
	craftService craftrouteservice.ICraftRouteService,
	alertRuleService alertservice.AlertLogService,
	deviceMonitorCalcDao devicemonitordao.DeviceMonitorCalcDao, cache cache.Cache) homeservice.HomeService {
	return &HomeServiceImpl{
		buildingDao:          buildingDao,
		deviceService:        deviceService,
		cache:                cache,
		craftService:         craftService,
		alertRuleService:     alertRuleService,
		deviceMonitorCalcDao: deviceMonitorCalcDao,
	}
}

func (h *HomeServiceImpl) GetHomeStats(c *gin.Context, isMobile bool) (*homemodels.HomeStats, error) {
	stats := &homemodels.HomeStats{}

	// 1. Get building count
	buildingCount, err := h.buildingDao.Count(c)
	if err != nil {
		zap.L().Error("build count error", zap.Error(err))
		return stats, err
	}
	stats.BuildingCount = buildingCount

	// 2. Get device stats
	statCount, err := h.deviceService.StatCount(c)
	if err != nil {
		zap.L().Error("device stat count error", zap.Error(err))
		return stats, err
	}
	if statCount == nil {
		return stats, errors.New("device stat count is nil")
	}
	stats.DeviceStats.TotalDevices = statCount.Total
	stats.DeviceStats.OnlineDevices = statCount.Online
	stats.DeviceStats.OfflineDevices = statCount.OffLine
	stats.DeviceStats.ExceptionCount = statCount.Exception
	stats.DeviceStats.MaintenanceCount = statCount.Maintenance

	// 计算调度策略的数量
	scheduleCount, err := h.craftService.Count(c)
	if err != nil {
		zap.L().Error("craftService count error", zap.Error(err))
		return stats, err
	}
	stats.ScheduleCount = scheduleCount

	// 计算告警策略的数量
	alertCount, err := h.alertRuleService.Count(c)
	if err != nil {
		zap.L().Error("alertRuleDao count error", zap.Error(err))
		return nil, err
	}
	stats.AlertsCount = alertCount

	// 实时告警监控
	alertList, err := h.alertRuleService.List(c, &alertmodels.SysAlertLogListReq{
		BaseEntityDQL: baize.BaseEntityDQL{
			Page: 1,
			Size: 5,
		},
	})
	if err != nil {
		zap.L().Error("alertRuleDao List error", zap.Error(err))
		return nil, err
	}
	if alertList != nil && alertList.Total > 0 {
		stats.Alerts = alertList.Rows
	}

	end := time.Now().UnixMilli()
	start := end - 86400*1000*6

	if !isMobile {
		// 读取主机性能
		monitorServer := monitormodels.NewServer()
		stats.Monitor = monitorServer

		endRankTime := time.Now().UnixMilli()
		startRankTime := end - 86400*1000
		// 读取最近10台的设备统计
		deviceRankCounter, err := h.deviceMonitorCalcDao.CounterByDevice(c, startRankTime, endRankTime, 10)
		if err != nil {
			zap.L().Error("device counter error align by device", zap.Error(err))
			return nil, err
		}
		stats.DeviceCounterRank = deviceRankCounter
	}

	// 读取每个小时的发送统计
	deviceCounter, err := h.deviceMonitorCalcDao.CounterByTimeRange(start, end, "1d")
	if err != nil {
		zap.L().Error("device counter error", zap.Error(err))
		return nil, err
	}
	stats.DeviceCounter = deviceCounter
	return stats, nil
}
