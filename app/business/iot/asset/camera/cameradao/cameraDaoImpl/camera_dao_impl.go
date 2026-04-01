package cameraDaoImpl

import (
	"context"
	"nova-factory-server/app/business/iot/asset/camera/cameradao"
	"nova-factory-server/app/business/iot/asset/camera/cameramodels"
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
func NewCameraDao(db *gorm.DB) cameradao.ICameraDao {
	return &CameraDaoImpl{
		db:    db,
		table: "iot_camera",
	}
}

// Create 创建摄像头
func (c *CameraDaoImpl) Create(ctx *gin.Context, camera *cameramodels.IotCamera) error {
	camera.Id = snowflake.GenID()
	camera.SetCreateBy(baizeContext.GetUserId(ctx))
	camera.DeptId = baizeContext.GetUserId(ctx)
	return c.db.Table(c.table).Create(camera).Error
}

// Update 更新摄像头
func (c *CameraDaoImpl) Update(ctx *gin.Context, camera *cameramodels.IotCamera) error {
	camera.SetUpdateBy(baizeContext.GetUserId(ctx))
	return c.db.Table(c.table).Save(camera).Error
}

// Delete 删除摄像头
func (c *CameraDaoImpl) Delete(ids []string) error {
	return c.db.Table(c.table).Where("id in (?)", ids).Delete(&cameramodels.IotCamera{}).Error
}

// GetById 根据ID获取摄像头
func (c *CameraDaoImpl) GetById(id int64) (*cameramodels.IotCamera, error) {
	var camera cameramodels.IotCamera
	err := c.db.Table(c.table).First(&camera, id).Error
	if err != nil {
		return nil, err
	}
	return &camera, nil
}

// GetByIds 根据多个ID批量获取摄像头
func (c *CameraDaoImpl) GetByIds(ids []int64) ([]*cameramodels.IotCamera, error) {
	var cameras []*cameramodels.IotCamera
	err := c.db.Table(c.table).Where("id IN ?", ids).Find(&cameras).Error
	if err != nil {
		return nil, err
	}
	return cameras, nil
}

// List 获取摄像头列表
func (c *CameraDaoImpl) List(req *cameramodels.IotCameraListReq) ([]*cameramodels.IotCamera, error) {
	var cameras []*cameramodels.IotCamera
	query := c.db.Table(c.table).Model(&cameramodels.IotCamera{})

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
func (c *CameraDaoImpl) Count(req *cameramodels.IotCameraListReq) (int64, error) {
	var count int64
	query := c.db.Table(c.table).Model(&cameramodels.IotCamera{})

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

// UpdateStatus 更新摄像头状态
func (c *CameraDaoImpl) UpdateStatus(ctx context.Context, id int64, status bool) error {
	return c.db.Table(c.table).Where("id = ?", id).Update("status", status).Error
}
