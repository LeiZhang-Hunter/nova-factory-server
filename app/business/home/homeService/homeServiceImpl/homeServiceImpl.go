package homeServiceImpl

import (
	"errors"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/business/asset/building/buildingDao"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorDao"
	"nova-factory-server/app/business/home/homeModels"
	"nova-factory-server/app/business/home/homeService"
	"nova-factory-server/app/business/monitor/monitorModels"
	"nova-factory-server/app/datasource/cache"
	"time"

	"github.com/gin-gonic/gin"
)

type HomeServiceImpl struct {
	buildingDao          buildingDao.BuildingDao
	deviceService        deviceService.IDeviceService
	cache                cache.Cache
	craftService         craftRouteService.ICraftRouteService
	alertRuleService     alertService.AlertLogService
	deviceMonitorCalcDao deviceMonitorDao.DeviceMonitorCalcDao
}

func NewHomeServiceImpl(buildingDao buildingDao.BuildingDao, deviceService deviceService.IDeviceService,
	craftService craftRouteService.ICraftRouteService,
	alertRuleService alertService.AlertLogService,
	deviceMonitorCalcDao deviceMonitorDao.DeviceMonitorCalcDao, cache cache.Cache) homeService.HomeService {
	return &HomeServiceImpl{
		buildingDao:          buildingDao,
		deviceService:        deviceService,
		cache:                cache,
		craftService:         craftService,
		alertRuleService:     alertRuleService,
		deviceMonitorCalcDao: deviceMonitorCalcDao,
	}
}

func (h *HomeServiceImpl) GetHomeStats(c *gin.Context) (*homeModels.HomeStats, error) {
	stats := &homeModels.HomeStats{}

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
	alertList, err := h.alertRuleService.List(c, &alertModels.SysAlertLogListReq{
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

	// 读取主机性能
	monitorServer := monitorModels.NewServer()
	stats.Monitor = monitorServer

	end := time.Now().UnixMilli()
	start := end - 86400*1000*6
	// 读取每个小时的发送统计
	deviceCounter, err := h.deviceMonitorCalcDao.CounterByTimeRange(start, end, "1d")
	if err != nil {
		zap.L().Error("device counter error", zap.Error(err))
		return nil, err
	}
	stats.DeviceCounter = deviceCounter

	endRankTime := time.Now().UnixMilli()
	startRankTime := end - 86400*1000*30
	// 读取最近10台的设备统计
	deviceRankCounter, err := h.deviceMonitorCalcDao.CounterByDevice(c, startRankTime, endRankTime, 10)
	if err != nil {
		zap.L().Error("device counter error align by device", zap.Error(err))
		return nil, err
	}
	stats.DeviceCounterRank = deviceRankCounter
	return stats, nil
}
