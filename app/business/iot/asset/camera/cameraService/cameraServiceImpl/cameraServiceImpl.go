package cameraServiceImpl

import (
	"nova-factory-server/app/business/iot/asset/camera/cameraDao"
	"nova-factory-server/app/business/iot/asset/camera/cameraModels"
	"nova-factory-server/app/business/iot/asset/camera/cameraService"
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
func (s *CameraServiceImpl) Create(camera *cameraModels.IotCamera) error {
	return s.cameraDao.Create(camera)
}

// Update 更新摄像头
func (s *CameraServiceImpl) Update(camera *cameraModels.IotCamera) error {
	return s.cameraDao.Update(camera)
}

// Delete 删除摄像头
func (s *CameraServiceImpl) Delete(id int64) error {
	return s.cameraDao.Delete(id)
}

// GetById 根据ID获取摄像头
func (s *CameraServiceImpl) GetById(id int64) (*cameraModels.IotCamera, error) {
	return s.cameraDao.GetById(id)
}

// List 获取摄像头列表
func (s *CameraServiceImpl) List() ([]*cameraModels.IotCamera, error) {
	return s.cameraDao.List()
}
