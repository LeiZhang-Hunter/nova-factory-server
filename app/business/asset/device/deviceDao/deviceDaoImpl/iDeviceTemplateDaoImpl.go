package deviceDaoImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
)

type IDeviceTemplateDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewIDeviceTemplateDaoImpl(ms *gorm.DB) deviceDao.IDeviceTemplateDao {
	return &IDeviceTemplateDaoImpl{
		db:        ms,
		tableName: "sys_device_template",
	}
}

func (i *IDeviceTemplateDaoImpl) GetById(c *gin.Context, id int64) (*deviceModels.SysDeviceTemplate, error) {
	var data *deviceModels.SysDeviceTemplate
	ret := i.db.Table(i.tableName).Where("template_id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&data)
	return data, ret.Error
}

func (i *IDeviceTemplateDaoImpl) Add(c *gin.Context, template *deviceModels.SysDeviceTemplate) (*deviceModels.SysDeviceTemplate, error) {
	ret := i.db.Table(i.tableName).Create(template)
	return template, ret.Error
}

func (i *IDeviceTemplateDaoImpl) Update(c *gin.Context, template *deviceModels.SysDeviceTemplate) (*deviceModels.SysDeviceTemplate, error) {
	ret := i.db.Table(i.tableName).Where("template_id = ?", template.TemplateID).Updates(template)
	return template, ret.Error
}

func (i *IDeviceTemplateDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.tableName).Where("template_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (i *IDeviceTemplateDaoImpl) List(c *gin.Context, req *deviceModels.SysDeviceTemplateDQL) (*deviceModels.SysDeviceTemplateListData, error) {
	db := i.db.Table(i.tableName)
	if req == nil {
		req = &deviceModels.SysDeviceTemplateDQL{}
	}
	if req.Protocol != "" {
		db = db.Where("protocol = ?", req.Protocol)
	}
	if req != nil && req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	size := 0
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}

	db = baizeContext.GetGormDataScope(c, db)
	db = db.Where("state = ?", 0)

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &deviceModels.SysDeviceTemplateListData{
			Rows:  make([]*deviceModels.SysDeviceTemplate, 0),
			Total: 0,
		}, ret.Error
	}
	var dto []*deviceModels.SysDeviceTemplate
	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &deviceModels.SysDeviceTemplateListData{
			Rows:  make([]*deviceModels.SysDeviceTemplate, 0),
			Total: 0,
		}, ret.Error
	}
	return &deviceModels.SysDeviceTemplateListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *IDeviceTemplateDaoImpl) GetByIds(c *gin.Context, ids []uint64) ([]*deviceModels.SysDeviceTemplate, error) {
	if ids == nil || len(ids) == 0 {
		return nil, errors.New("ids is null")
	}
	var dto []*deviceModels.SysDeviceTemplate
	ret := i.db.Table(i.tableName).Where("template_id in (?)", ids).Where("state = ?", commonStatus.NORMAL).Find(&dto)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return dto, nil

}
