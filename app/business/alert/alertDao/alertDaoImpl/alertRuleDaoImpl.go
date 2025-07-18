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

type AlertRuleDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAlertRuleDaoImpl(db *gorm.DB) alertDao.AlertRuleDao {
	return &AlertRuleDaoImpl{
		db:    db,
		table: "sys_alert",
	}
}

func (ac *AlertRuleDaoImpl) Create(c *gin.Context, data *alertModels.SetSysAlert) (*alertModels.SysAlert, error) {
	data.ID = snowflake.GenID()
	value := alertModels.ToSysAlert(data)
	value.SetCreateBy(baizeContext.GetUserId(c))
	value.DeptID = baizeContext.GetDeptId(c)
	ret := ac.db.Table(ac.table).Create(&value)
	return value, ret.Error
}

func (ac *AlertRuleDaoImpl) Update(c *gin.Context, data *alertModels.SetSysAlert) (*alertModels.SysAlert, error) {
	value := alertModels.ToSysAlert(data)
	value.SetUpdateBy(baizeContext.GetUserId(c))
	ret := ac.db.Table(ac.table).Where("id = ?", value.ID).Updates(&value)
	return value, ret.Error
}

func (ac *AlertRuleDaoImpl) List(c *gin.Context, req *alertModels.SysAlertListReq) (*alertModels.SysAlertList, error) {
	db := ac.db.Table(ac.table)

	if req != nil && req.GatewayID != 0 {
		db = db.Where("gateway_id = ?", req.GatewayID)
	}
	if req != nil && req.TemplateID != 0 {
		db = db.Where("template_id = ?", req.TemplateID)
	}
	if req != nil && req.Name != "" {
		db = db.Where("name like %s", "%"+req.Name+"%")
	}
	if req != nil && req.Status != nil {
		db = db.Where("status = ?", req.Status)
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
	var dto []*alertModels.SysAlert

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &alertModels.SysAlertList{
			Rows:  make([]*alertModels.SetSysAlert, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &alertModels.SysAlertList{
			Rows:  make([]*alertModels.SetSysAlert, 0),
			Total: 0,
		}, ret.Error
	}

	list := make([]*alertModels.SetSysAlert, 0)
	for _, value := range dto {
		list = append(list, alertModels.FromSysAlertToSetData(value))
	}
	return &alertModels.SysAlertList{
		Rows:  list,
		Total: uint64(total),
	}, nil
}

func (ac *AlertRuleDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := ac.db.Table(ac.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (ac *AlertRuleDaoImpl) GetByGatewayId(c *gin.Context, gatewayId uint64) (*alertModels.SysAlert, error) {
	var dto *alertModels.SysAlert
	ret := ac.db.Table(ac.table).Where("gateway_id = ?", gatewayId).Where("state = ?", commonStatus.NORMAL).First(&dto)
	return dto, ret.Error
}

func (ac *AlertRuleDaoImpl) GetById(c *gin.Context, id uint64) (*alertModels.SysAlert, error) {
	var dto *alertModels.SysAlert
	ret := ac.db.Table(ac.table).Where("id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&dto)
	return dto, ret.Error
}

func (ac *AlertRuleDaoImpl) Change(c *gin.Context, data *alertModels.ChangeSysAlert) error {
	ret := ac.db.Table(ac.table).Where("id = ?", data.ID).Update("status", data.Status)
	return ret.Error
}
