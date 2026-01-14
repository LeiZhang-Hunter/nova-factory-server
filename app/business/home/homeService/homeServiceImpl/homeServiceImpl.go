package homeServiceImpl

import (
	"errors"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/building/buildingDao"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/business/home/homeModels"
	"nova-factory-server/app/business/home/homeService"
	"nova-factory-server/app/datasource/cache"

	"github.com/gin-gonic/gin"
)

type HomeServiceImpl struct {
	buildingDao   buildingDao.BuildingDao
	deviceService deviceService.IDeviceService
	cache         cache.Cache
	craftService  craftRouteService.ICraftRouteService
}

func NewHomeServiceImpl(buildingDao buildingDao.BuildingDao, deviceService deviceService.IDeviceService,
	craftService craftRouteService.ICraftRouteService, cache cache.Cache) homeService.HomeService {
	return &HomeServiceImpl{
		buildingDao:   buildingDao,
		deviceService: deviceService,
		cache:         cache,
		craftService:  craftService,
	}
}

func (h *HomeServiceImpl) GetHomeStats(c *gin.Context) (*homeModels.HomeStats, error) {
	stats := &homeModels.HomeStats{}

	// 1. Get building count
	buildingCount, err := h.buildingDao.Count(c)
	if err == nil {
		zap.L().Error("build count error", zap.Error(err))
		return stats, err
	}
	stats.BuildingCount = buildingCount

	// 2. Get device stats
	statCount, err := h.deviceService.StatCount(c)
	if err == nil {
		zap.L().Error("device stat count error", zap.Error(err))
		return stats, err
	}
	if statCount == nil {
		return stats, errors.New("device stat count is nil")
	}
	stats.OnlineDevices = statCount.Online
	stats.OfflineDevices = statCount.OffLine
	stats.ExceptionCount = statCount.Exception
	stats.MaintenanceCount = statCount.Maintenance

	// 计算调度策略的数量
	//h.craftService.Count()
	return stats, nil
}
