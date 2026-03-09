package cameraDaoImpl

import (
	"nova-factory-server/app/business/iot/asset/camera/cameraDao"
	"nova-factory-server/app/business/iot/asset/camera/cameraModels"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CameraDaoImpl 摄像头数据访问实现
type CameraDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewCameraDao 创建摄像头数据访问实例
func NewCameraDao(db *gorm.DB) cameraDao.ICameraDao {
	return &CameraDaoImpl{
		db:    db,
		table: "iot_camera",
	}
}

// Create 创建摄像头
func (c *CameraDaoImpl) Create(ctx *gin.Context, camera *cameraModels.IotCamera) error {
	camera.Id = snowflake.GenID()
	camera.SetCreateBy(baizeContext.GetUserId(ctx))
	camera.DeptId = baizeContext.GetUserId(ctx)
	return c.db.Table(c.table).Create(camera).Error
}

// Update 更新摄像头
func (c *CameraDaoImpl) Update(ctx *gin.Context, camera *cameraModels.IotCamera) error {
	camera.SetUpdateBy(baizeContext.GetUserId(ctx))
	return c.db.Table(c.table).Save(camera).Error
}

// Delete 删除摄像头
func (c *CameraDaoImpl) Delete(ids []string) error {
	return c.db.Table(c.table).Where("id in (?)", ids).Delete(&cameraModels.IotCamera{}).Error
}

// GetById 根据ID获取摄像头
func (c *CameraDaoImpl) GetById(id int64) (*cameraModels.IotCamera, error) {
	var camera cameraModels.IotCamera
	err := c.db.Table(c.table).First(&camera, id).Error
	if err != nil {
		return nil, err
	}
	return &camera, nil
}

// List 获取摄像头列表
func (c *CameraDaoImpl) List(req *cameraModels.IotCameraListReq) ([]*cameraModels.IotCamera, error) {
	var cameras []*cameraModels.IotCamera
	query := c.db.Table(c.table).Model(&cameraModels.IotCamera{})

	if req != nil {
		if req.Name != "" {
			query = query.Where("name LIKE ?", "%"+req.Name+"%")
		}
		if req.IpAddress != "" {
			query = query.Where("ip_address LIKE ?", "%"+req.IpAddress+"%")
		}
		if req.Brand != "" {
			query = query.Where("brand = ?", req.Brand)
		}
	}

	err := query.Find(&cameras).Error
	if err != nil {
		return nil, err
	}
	return cameras, nil
}

// Count 获取摄像头总数
func (c *CameraDaoImpl) Count(req *cameraModels.IotCameraListReq) (int64, error) {
	var count int64
	query := c.db.Table(c.table).Model(&cameraModels.IotCamera{})

	if req != nil {
		if req.Name != "" {
			query = query.Where("name LIKE ?", "%"+req.Name+"%")
		}
		if req.IpAddress != "" {
			query = query.Where("ip_address LIKE ?", "%"+req.IpAddress+"%")
		}
		if req.Brand != "" {
			query = query.Where("brand LIKE ?", "%"+req.Brand+"%")
		}
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
