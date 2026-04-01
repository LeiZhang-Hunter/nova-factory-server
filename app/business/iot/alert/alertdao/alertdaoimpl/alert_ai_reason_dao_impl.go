package alertdaoimpl

import (
	"errors"
	"nova-factory-server/app/business/iot/alert/alertdao"
	"nova-factory-server/app/business/iot/alert/alertmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AlertAiReasonDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAlertAiReasonDaoImpl(db *gorm.DB) alertdao.AlertAiReasonDao {
	return &AlertAiReasonDaoImpl{
		db:    db,
		table: "sys_alert_ai_reason",
	}
}

func (a *AlertAiReasonDaoImpl) Set(c *gin.Context, data *alertmodels.SetAlertAiReason) (*alertmodels.SysAlertAiReason, error) {
	value := alertmodels.FromSetAlertReasonToData(data)
	if value.ID == 0 {
		value.SetCreateBy(baizeContext.GetUserId(c))
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		ret := a.db.Table(a.table).Create(&value)
		return value, ret.Error
	} else {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		ret := a.db.Table(a.table).Debug().Where("id = ?", value.ID).Updates(&value)
		return value, ret.Error
	}
}

func (a *AlertAiReasonDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := a.db.Table(a.table).Where("id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (a *AlertAiReasonDaoImpl) List(c *gin.Context, req *alertmodels.SysAlertAiReasonReq) (*alertmodels.SysAlertReasonList, error) {
	db := a.db.Table(a.table)

	if req != nil && req.Name != "" {
		db = db.Where("name like ?", "%"+req.Name+"%")
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
	var dto []*alertmodels.SysAlertAiReason

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &alertmodels.SysAlertReasonList{
			Rows:  make([]*alertmodels.SysAlertAiReason, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &alertmodels.SysAlertReasonList{
			Rows:  make([]*alertmodels.SysAlertAiReason, 0),
			Total: 0,
		}, ret.Error
	}

	return &alertmodels.SysAlertReasonList{
		Rows:  dto,
		Total: uint64(total),
	}, nil
}

func (a *AlertAiReasonDaoImpl) GetById(c *gin.Context, id int64) (*alertmodels.SysAlertAiReason, error) {
	var dto *alertmodels.SysAlertAiReason
	ret := a.db.Table(a.table).Where("id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&dto)
	if ret.Error != nil {
		if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return dto, ret.Error
}
