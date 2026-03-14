package cameraService

import (
	"nova-factory-server/app/business/iot/asset/camera/cameraModels"

	"github.com/gin-gonic/gin"
)

// ICameraService 摄像头服务接口
type ICameraService interface {
	// Create 创建摄像头
	Create(ctx *gin.Context, camera *cameraModels.IotCamera) error
	// Update 更新摄像头
	Update(ctx *gin.Context, camera *cameraModels.IotCamera) error
	// Delete 删除摄像头
	Delete(id []string) error
	// GetById 根据ID获取摄像头
	GetById(id int64) (*cameraModels.IotCamera, error)
	// GetDetailById 根据ID获取摄像头并合并实时信息
	GetDetailById(ctx *gin.Context, id int64) (*cameraModels.IotCameraDetail, error)
	// List 获取摄像头列表
	List(req *cameraModels.IotCameraListReq) (*cameraModels.IotCameraList, error)
}
