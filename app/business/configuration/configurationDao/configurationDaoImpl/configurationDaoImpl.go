package configurationDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/configuration/configurationDao"
	"nova-factory-server/app/business/configuration/configurationModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type ConfigurationDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewConfigurationDaoImpl(db *gorm.DB) configurationDao.ConfigurationDao {
	return &ConfigurationDaoImpl{
		db:    db,
		table: "sys_configuration",
	}
}

func (i *ConfigurationDaoImpl) List(c *gin.Context, req *configurationModels.SysConfigurationReq) (*configurationModels.SysConfigurationList, error) {
	db := i.db.Table(i.table)

	if req != nil && req.Name != "" {
		db = db.Where("name like ?", "%"+req.Name+"%")
	}
	if req != nil && req.Tag != "" {
		db = db.Where("tag = ?", req.Tag)
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

	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)

	var dto []*configurationModels.SysConfiguration
	var total int64

	ret := db.Count(&total)
	if ret.Error != nil {
		return &configurationModels.SysConfigurationList{
			Rows:  make([]*configurationModels.SysConfiguration, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Debug().Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &configurationModels.SysConfigurationList{
			Rows:  make([]*configurationModels.SysConfiguration, 0),
			Total: 0,
		}, ret.Error
	}

	return &configurationModels.SysConfigurationList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *ConfigurationDaoImpl) Set(c *gin.Context, data *configurationModels.SetSysConfiguration) (*configurationModels.SysConfiguration, error) {
	value := configurationModels.ToSysConfiguration(data)
	if data.ID == 0 {
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		value.SetCreateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Create(value)
		return value, ret.Error
	}
	value.SetUpdateBy(baizeContext.GetUserId(c))
	ret := i.db.Table(i.table).Where("id = ?", data.ID).Updates(value)
	return value, ret.Error
}

func (i *ConfigurationDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.table).Where("id in (?)", ids).Delete(&configurationModels.SysConfiguration{})
	return ret.Error
}
