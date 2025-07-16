package alertDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type AlertSinkTemplateDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAlertSinkTemplateDaoImpl(db *gorm.DB) alertDao.AlertSinkTemplateDao {
	return &AlertSinkTemplateDaoImpl{
		db:    db,
		table: "sys_alert_sink_template",
	}
}

func (ac *AlertSinkTemplateDaoImpl) Create(c *gin.Context, data *alertModels.SetSysAlertSinkTemplate) (*alertModels.SysAlertSinkTemplate, error) {
	data.ID = snowflake.GenID()
	value := alertModels.ToSysAlertSinkTemplate(data)
	value.SetCreateBy(baizeContext.GetUserId(c))
	value.DeptID = baizeContext.GetDeptId(c)
	ret := ac.db.Table(ac.table).Create(&value)
	return value, ret.Error
}

func (ac *AlertSinkTemplateDaoImpl) Update(c *gin.Context, data *alertModels.SetSysAlertSinkTemplate) (*alertModels.SysAlertSinkTemplate, error) {
	value := alertModels.ToSysAlertSinkTemplate(data)
	value.SetUpdateBy(baizeContext.GetUserId(c))
	ret := ac.db.Table(ac.table).Where("id = ?", value.ID).Updates(&value)
	return value, ret.Error
}

func (ac *AlertSinkTemplateDaoImpl) List(c *gin.Context, req *alertModels.SysAlertSinkTemplateReq) (*alertModels.SysAlertSinkTemplateListData, error) {
	db := ac.db.Table(ac.table)

	if req != nil && req.GatewayID != 0 {
		db = db.Where("gateway_id = ?", req.GatewayID)
	}

	if req != nil && req.Name != "" {
		db = db.Where("name like %s", "%"+req.Name+"%")
	}
	if req != nil && req.Addr != "" {
		db = db.Where("addr like ?", "%"+req.Addr+"%")
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
	var dto []*alertModels.SysAlertSinkTemplate

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &alertModels.SysAlertSinkTemplateListData{
			Rows:  make([]*alertModels.SysAlertSinkTemplate, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &alertModels.SysAlertSinkTemplateListData{
			Rows:  make([]*alertModels.SysAlertSinkTemplate, 0),
			Total: 0,
		}, ret.Error
	}
	return &alertModels.SysAlertSinkTemplateListData{
		Rows:  dto,
		Total: uint64(total),
	}, nil
}

func (ac *AlertSinkTemplateDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := ac.db.Table(ac.table).Where("id in (?)", ids).Delete(&alertModels.SysAlertSinkTemplate{})
	return ret.Error
}

func (ac *AlertSinkTemplateDaoImpl) GetByGatewayId(c *gin.Context, gatewayId uint64) (*alertModels.SysAlertSinkTemplate, error) {
	var dto *alertModels.SysAlertSinkTemplate
	ret := ac.db.Table(ac.table).Where("gateway_id = ?", gatewayId).Where("state = ?", commonStatus.NORMAL).First(&dto)
	return dto, ret.Error
}

func (ac *AlertSinkTemplateDaoImpl) GetById(c *gin.Context, id uint64) (*alertModels.SysAlertSinkTemplate, error) {
	var dto *alertModels.SysAlertSinkTemplate
	ret := ac.db.Table(ac.table).Where("id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&dto)
	return dto, ret.Error
}
