package cameraServiceImpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"nova-factory-server/app/business/iot/asset/camera/cameraDao"
	"nova-factory-server/app/business/iot/asset/camera/cameraModels"
	"nova-factory-server/app/business/iot/asset/camera/cameraService"
	redisConst "nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/datasource/cache/cacheError"

	"github.com/gin-gonic/gin"
	v1 "github.com/novawatcher-io/nova-factory-payload/camera/v1"
	"go.uber.org/zap"
)

// CameraServiceImpl 摄像头服务实现
type CameraServiceImpl struct {
	cameraDao cameraDao.ICameraDao
}

// NewCameraService 创建摄像头服务实例
func NewCameraService(cameraDao cameraDao.ICameraDao) cameraService.ICameraService {
	return &CameraServiceImpl{cameraDao: cameraDao}
}

// Create 创建摄像头
func (s *CameraServiceImpl) Create(ctx *gin.Context, camera *cameraModels.IotCamera) error {
	return s.cameraDao.Create(ctx, camera)
}

// Update 更新摄像头
func (s *CameraServiceImpl) Update(ctx *gin.Context, camera *cameraModels.IotCamera) error {
	return s.cameraDao.Update(ctx, camera)
}

// Delete 删除摄像头
func (s *CameraServiceImpl) Delete(idList []string) error {
	return s.cameraDao.Delete(idList)
}

// GetById 根据ID获取摄像头
func (s *CameraServiceImpl) GetById(id int64) (*cameraModels.IotCamera, error) {
	return s.cameraDao.GetById(id)
}

// GetDetailById 根据ID获取摄像头并合并实时信息
func (s *CameraServiceImpl) GetDetailById(ctx *gin.Context, id int64) (*cameraModels.IotCameraDetail, error) {
	cameraInfo, err := s.cameraDao.GetById(id)
	if err != nil {
		return nil, err
	}
	detail := &cameraModels.IotCameraDetail{
		IotCamera: *cameraInfo,
	}

	cacheClient := cache.NewCache()
	cacheValue, cacheErr := cacheClient.Get(ctx, fmt.Sprintf(redisConst.CameraInfoCacheKey, id))
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
		}
	} else if !errors.Is(cacheErr, cacheError.Nil) {
		zap.L().Warn("read camera realtime cache error", zap.Error(cacheErr), zap.Int64("camera_id", id))
	}
	return detail, nil
}

// List 获取摄像头列表
func (s *CameraServiceImpl) List(req *cameraModels.IotCameraListReq) (*cameraModels.IotCameraList, error) {
	cameras, err := s.cameraDao.List(req)
	if err != nil {
		return nil, err
	}

	total, err := s.cameraDao.Count(req)
	if err != nil {
		return nil, err
	}

	return &cameraModels.IotCameraList{
		Rows:  cameras,
		Total: total,
	}, nil
}
