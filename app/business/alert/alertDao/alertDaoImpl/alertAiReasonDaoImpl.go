package alertDaoImpl

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type AlertAiReasonDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAlertAiReasonDaoImpl(db *gorm.DB) alertDao.AlertAiReasonDao {
	return &AlertAiReasonDaoImpl{
		db:    db,
		table: "sys_alert_ai_reason",
	}
}

func (a *AlertAiReasonDaoImpl) Set(c *gin.Context, data *alertModels.SetAlertAiReason) (*alertModels.SysAlertAiReason, error) {
	value := alertModels.FromSetAlertReasonToData(data)
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

func (a *AlertAiReasonDaoImpl) List(c *gin.Context, req *alertModels.SysAlertAiReasonReq) (*alertModels.SysAlertReasonList, error) {
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
	var dto []*alertModels.SysAlertAiReason

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &alertModels.SysAlertReasonList{
			Rows:  make([]*alertModels.SysAlertAiReason, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &alertModels.SysAlertReasonList{
			Rows:  make([]*alertModels.SysAlertAiReason, 0),
			Total: 0,
		}, ret.Error
	}

	for k, v := range dto {
		dto[k].DatasetIdList = make([]string, 0)
		err := json.Unmarshal([]byte(v.DatasetIds), &dto[k].DatasetIdList)
		if err != nil {
			zap.L().Error("json unmarshal error", zap.Error(err))
		}
	}

	return &alertModels.SysAlertReasonList{
		Rows:  dto,
		Total: uint64(total),
	}, nil
}
