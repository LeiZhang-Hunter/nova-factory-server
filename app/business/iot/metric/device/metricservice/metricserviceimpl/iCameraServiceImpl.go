package metricserviceimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"nova-factory-server/app/business/iot/asset/camera/cameradao"
	"nova-factory-server/app/business/iot/asset/camera/cameramodels"
	"nova-factory-server/app/business/iot/metric/device/metricservice"
	"nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
	"time"

	v1 "github.com/novawatcher-io/nova-factory-payload/camera/v1"
	"go.uber.org/zap"
)

type ICameraServiceImpl struct {
	cache     cache.Cache
	cameraDao cameradao.ICameraDao
}

func NewICameraServiceImpl(cache cache.Cache, cameraDao cameradao.ICameraDao) metricservice.ICameraService {
	return &ICameraServiceImpl{
		cache:     cache,
		cameraDao: cameraDao,
	}
}

// Report 接收数据上报，存储摄像头信息
func (c *ICameraServiceImpl) Report(ctx context.Context, req *v1.CameraData) error {
	latestReports := make(map[uint64]*v1.CameraRequest)
	for _, report := range req.Request {
		id := report.GetDeviceId()
		if existing, ok := latestReports[id]; !ok || report.GetTime() > existing.GetTime() {
			latestReports[id] = report
		}
	}

	for _, report := range latestReports {
		// Save camera info/data to Redis
		data, err := json.Marshal(report)
		if err != nil {
			zap.L().Error("Report error", zap.Error(err))
			continue
		}
		infoKey := fmt.Sprintf(redis.CameraInfoCacheKey, report.GetDeviceId())
		c.cache.Set(ctx, infoKey, string(data), 2*time.Minute)
	}

	// Optimization: Batch read camera DB info
	deviceIds := make([]int64, 0, len(latestReports))
	for id := range latestReports {
		deviceIds = append(deviceIds, int64(id))
	}

	cameras, err := c.cameraDao.GetByIds(deviceIds)
	if err != nil {
		zap.L().Error("Batch get cameras from DB error", zap.Error(err))
		// Log error but maybe continue with some best-effort?
		// Actually, if DB fails, we can't do the status update check.
		return err
	}

	cameraMap := make(map[int64]*cameramodels.IotCamera)
	for _, cam := range cameras {
		cameraMap[cam.Id] = cam
	}

	for id, report := range latestReports {
		deviceId := int64(id)
		camera, ok := cameraMap[deviceId]
		if !ok {
			continue
		}

		// Optimization: Read camera DB info, if status and isOnLine are the same, do nothing.
		// Otherwise, update status in DB.
		isOnLine := report.GetIsOnLine()
		if camera.Status == nil || *camera.Status != isOnLine {
			err = c.cameraDao.UpdateStatus(ctx, deviceId, isOnLine)
			if err != nil {
				zap.L().Error("Update camera status error", zap.Error(err), zap.Int64("deviceId", deviceId))
			} else {
				zap.L().Debug("Camera status updated", zap.Int64("deviceId", deviceId), zap.Bool("isOnLine", isOnLine))
			}
		}
	}
	return nil
}
