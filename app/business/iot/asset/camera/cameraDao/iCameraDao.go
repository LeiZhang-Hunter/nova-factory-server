package cameraDao

import (
	"nova-factory-server/app/business/iot/asset/camera/cameraModels"
)

// ICameraDao 摄像头数据访问接口
type ICameraDao interface {
	// Create 创建摄像头
	Create(camera *cameraModels.IotCamera) error
	// Update 更新摄像头
	Update(camera *cameraModels.IotCamera) error
	// Delete 删除摄像头
	Delete(ids []string) error
	// GetById 根据ID获取摄像头
	GetById(id int64) (*cameraModels.IotCamera, error)
	// List 获取摄像头列表
	List(req *cameraModels.IotCameraListReq) ([]*cameraModels.IotCamera, error)
	// Count 获取摄像头总数
	Count(req *cameraModels.IotCameraListReq) (int64, error)
}
