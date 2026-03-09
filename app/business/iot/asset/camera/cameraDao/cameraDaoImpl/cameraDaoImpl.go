package cameraDaoImpl

import (
	"nova-factory-server/app/business/iot/asset/camera/cameraDao"
	"nova-factory-server/app/business/iot/asset/camera/cameraModels"

	"gorm.io/gorm"
)

// CameraDaoImpl 摄像头数据访问实现
type CameraDaoImpl struct {
	db *gorm.DB
}

// NewCameraDao 创建摄像头数据访问实例
func NewCameraDao(db *gorm.DB) cameraDao.ICameraDao {
	return &CameraDaoImpl{db: db}
}

// Create 创建摄像头
func (c *CameraDaoImpl) Create(camera *cameraModels.IotCamera) error {
	return c.db.Create(camera).Error
}

// Update 更新摄像头
func (c *CameraDaoImpl) Update(camera *cameraModels.IotCamera) error {
	return c.db.Save(camera).Error
}

// Delete 删除摄像头
func (c *CameraDaoImpl) Delete(id int64) error {
	return c.db.Delete(&cameraModels.IotCamera{}, id).Error
}

// GetById 根据ID获取摄像头
func (c *CameraDaoImpl) GetById(id int64) (*cameraModels.IotCamera, error) {
	var camera cameraModels.IotCamera
	err := c.db.First(&camera, id).Error
	if err != nil {
		return nil, err
	}
	return &camera, nil
}

// List 获取摄像头列表
func (c *CameraDaoImpl) List() ([]*cameraModels.IotCamera, error) {
	var cameras []*cameraModels.IotCamera
	err := c.db.Find(&cameras).Error
	if err != nil {
		return nil, err
	}
	return cameras, nil
}
