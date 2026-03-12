package metricServiceImpl

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/novawatcher-io/nova-factory-payload/camera/v1"
	"go.uber.org/zap"
	"nova-factory-server/app/business/iot/metric/device/metricService"
	"nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
	"time"
)

type ICameraServiceImpl struct {
	cache cache.Cache
}

func NewICameraServiceImpl(cache cache.Cache) metricService.ICameraService {
	return &ICameraServiceImpl{
		cache: cache,
	}
}

func (c *ICameraServiceImpl) Report(ctx context.Context, req *v1.CameraData) error {
	for _, report := range req.Request {
		// Save camera info/data to Redis
		data, err := json.Marshal(report)
		if err != nil {
			zap.L().Error("Report error", zap.Error(err))
			continue
		}
		infoKey := fmt.Sprintf(redis.CameraInfoCacheKey, report.DeviceId)
		c.cache.Set(ctx, infoKey, string(data), 2*time.Minute)

	}
	return nil
}
