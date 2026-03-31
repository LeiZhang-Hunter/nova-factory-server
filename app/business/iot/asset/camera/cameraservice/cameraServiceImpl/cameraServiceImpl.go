package cameraServiceImpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"nova-factory-server/app/business/iot/asset/camera/cameradao"
	"nova-factory-server/app/business/iot/asset/camera/cameramodels"
	"nova-factory-server/app/business/iot/asset/camera/cameraservice"
	redisConst "nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/datasource/cache/cacheError"

	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/camera/v1"
	"go.uber.org/zap"
)

// CameraServiceImpl 摄像头服务实现
type CameraServiceImpl struct {
	cameraDao cameradao.ICameraDao
	cache     cache.Cache
}

// NewCameraService 创建摄像头服务实例
func NewCameraService(cameraDao cameradao.ICameraDao, cache cache.Cache) cameraservice.ICameraService {
	return &CameraServiceImpl{cameraDao: cameraDao, cache: cache}
}

// Create 创建摄像头
func (s *CameraServiceImpl) Create(ctx *gin.Context, camera *cameramodels.IotCamera) error {
	return s.cameraDao.Create(ctx, camera)
}

// Update 更新摄像头
func (s *CameraServiceImpl) Update(ctx *gin.Context, camera *cameramodels.IotCamera) error {
	return s.cameraDao.Update(ctx, camera)
}

// Delete 删除摄像头
func (s *CameraServiceImpl) Delete(idList []string) error {
	return s.cameraDao.Delete(idList)
}

// GetById 根据ID获取摄像头
func (s *CameraServiceImpl) GetById(id int64) (*cameramodels.IotCamera, error) {
	return s.cameraDao.GetById(id)
}

// GetDetailById 根据ID获取摄像头并合并实时信息
func (s *CameraServiceImpl) GetDetailById(ctx *gin.Context, id int64) (*cameramodels.IotCameraDetail, error) {
	cameraInfo, err := s.cameraDao.GetById(id)
	if err != nil {
		return nil, err
	}
	detail := &cameramodels.IotCameraDetail{
		IotCamera: *cameraInfo,
	}

	cacheValue, cacheErr := s.cache.Get(ctx, fmt.Sprintf(redisConst.CameraInfoCacheKey, id))
	if cacheErr == nil {
		realtimeData := new(v1.CameraRequest)
		if err = json.Unmarshal([]byte(cacheValue), realtimeData); err == nil {
			status := realtimeData.GetIsOnLine()
			detail.Status = &status
			realtime := make(map[string]interface{})
			cacheBytes, marshalErr := json.Marshal(realtimeData)
			if marshalErr == nil {
				_ = json.Unmarshal(cacheBytes, &realtime)
			}
			detail.GatewayRealtime = realtime

			// 检查状态是否不一致，不一致则同步更新
			if cameraInfo.Status == nil || *cameraInfo.Status != status {
				// 更新数据库中的状态
				cameraInfo.Status = &status
				if err := s.cameraDao.Update(ctx, cameraInfo); err != nil {
					zap.L().Warn("update camera status error", zap.Error(err), zap.Int64("camera_id", id))
				}
			}
		}
	} else if !errors.Is(cacheErr, cacheError.Nil) {
		zap.L().Warn("read camera realtime cache error", zap.Error(cacheErr), zap.Int64("camera_id", id))
	}
	return detail, nil
}

// List 获取摄像头列表
func (s *CameraServiceImpl) List(c *gin.Context, req *cameramodels.IotCameraListReq) (*cameramodels.IotCameraList, error) {
	cameras, err := s.cameraDao.List(req)
	if err != nil {
		return nil, err
	}

	total, err := s.cameraDao.Count(req)
	if err != nil {
		return nil, err
	}

	// 从 Redis 批量读取实时状态并更新
	if len(cameras) > 0 {

		// 构建所有摄像头的缓存键
		keys := make([]string, len(cameras))
		for i, camera := range cameras {
			keys[i] = fmt.Sprintf(redisConst.CameraInfoCacheKey, camera.Id)
		}

		// 批量读取 Redis 缓存
		slice := s.cache.MGet(c, keys)
		vals := slice.Val()

		// 用于存储需要更新状态的摄像头
		camerasToUpdate := make([]*cameramodels.IotCamera, 0)

		// 处理缓存结果
		for i, val := range vals {
			camera := cameras[i]

			// 确定在线状态
			isOnline := false
			if val != nil {
				// 尝试将 val 转换为字符串
				str, ok := val.(string)
				if ok && str != "" {
					realtimeData := new(v1.CameraRequest)
					if err = json.Unmarshal([]byte(str), realtimeData); err == nil {
						// 使用解码后的数据判断在线状态
						isOnline = realtimeData.GetIsOnLine()
					} else {
						// 解码失败，视为离线
						isOnline = false
					}
				} else {
					// 缓存值不是字符串，视为离线
					isOnline = false
				}
			} else {
				// 缓存中没有数据，视为离线
				isOnline = false
			}
			onlineBool := isOnline

			// 检查状态是否不一致
			statusChanged := false
			if camera.Status == nil {
				// 数据库中没有状态，设置为当前状态
				camera.Status = &onlineBool
				statusChanged = true
			} else if *camera.Status != isOnline {
				// 状态不一致，更新为当前状态
				camera.Status = &onlineBool
				statusChanged = true
			}

			// 如果状态发生变化，添加到更新列表
			if statusChanged {
				camerasToUpdate = append(camerasToUpdate, camera)
			}
		}

		// 批量更新状态（如果有需要更新的摄像头）
		for _, camera := range camerasToUpdate {
			if err := s.cameraDao.Update(c, camera); err != nil {
				zap.L().Warn("update camera status error", zap.Error(err), zap.Int64("camera_id", camera.Id))
			}
		}
	}

	return &cameramodels.IotCameraList{
		Rows:  cameras,
		Total: total,
	}, nil
}
