package cameraService

import (
	"nova-factory-server/app/business/iot/asset/camera/cameraModels"
)

// ICameraService 摄像头服务接口
type ICameraService interface {
	// Create 创建摄像头
	Create(camera *cameraModels.IotCamera) error
	// Update 更新摄像头
	Update(camera *cameraModels.IotCamera) error
	// Delete 删除摄像头
	Delete(id []string) error
	// GetById 根据ID获取摄像头
	GetById(id int64) (*cameraModels.IotCamera, error)
	// List 获取摄像头列表
	List(req *cameraModels.IotCameraListReq) (*cameraModels.IotCameraList, error)
}
