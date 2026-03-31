package configurationDaoImpl

import (
	"nova-factory-server/app/business/iot/configuration/configurationdao"
	"nova-factory-server/app/business/iot/configuration/configurationmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConfigurationDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewConfigurationDaoImpl(db *gorm.DB) configurationdao.ConfigurationDao {
	return &ConfigurationDaoImpl{
		db:    db,
		table: "sys_configuration",
	}
}

func (i *ConfigurationDaoImpl) List(c *gin.Context, req *configurationmodels.SysConfigurationReq) (*configurationmodels.SysConfigurationList, error) {
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

	var dto []*configurationmodels.SysConfiguration
	var total int64

	ret := db.Count(&total)
	if ret.Error != nil {
		return &configurationmodels.SysConfigurationList{
			Rows:  make([]*configurationmodels.SysConfiguration, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Debug().Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &configurationmodels.SysConfigurationList{
			Rows:  make([]*configurationmodels.SysConfiguration, 0),
			Total: 0,
		}, ret.Error
	}

	return &configurationmodels.SysConfigurationList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (i *ConfigurationDaoImpl) Set(c *gin.Context, data *configurationmodels.SetSysConfiguration) (*configurationmodels.SysConfiguration, error) {
	value := configurationmodels.ToSysConfiguration(data)
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
	ret := i.db.Table(i.table).Where("id in (?)", ids).Delete(&configurationmodels.SysConfiguration{})
	return ret.Error
}
