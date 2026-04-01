package alertdaoimpl

import (
	"nova-factory-server/app/business/iot/alert/alertdao"
	"nova-factory-server/app/business/iot/alert/alertmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AlertLogDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAlertLogDaoImpl(db *gorm.DB) alertdao.AlertLogDao {
	return &AlertLogDaoImpl{
		db:    db,
		table: "sys_alert_log",
	}
}

func (log *AlertLogDaoImpl) Export(c *gin.Context, data []*alertmodels.SysAlertLog) error {
	ret := log.db.Table(log.table).Save(data)
	return ret.Error
}

func (log *AlertLogDaoImpl) List(c *gin.Context, req *alertmodels.SysAlertLogListReq) (*alertmodels.SysAlertLogList, error) {
	db := log.db.Table(log.table)

	if req != nil && req.GatewayID != 0 {
		db = db.Where("gateway_id = ?", req.GatewayID)
	}
	if req != nil && req.AlertID != 0 {
		db = db.Where("alert_id = ?", req.AlertID)
	}
	if req != nil && req.Message != "" {
		db = db.Where("message like %s", "%"+req.Message+"%")
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
	var dto []*alertmodels.SysAlertLog

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &alertmodels.SysAlertLogList{
			Rows:  make([]*alertmodels.SysAlertLog, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &alertmodels.SysAlertLogList{
			Rows:  make([]*alertmodels.SysAlertLog, 0),
			Total: 0,
		}, ret.Error
	}

	return &alertmodels.SysAlertLogList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (log *AlertLogDaoImpl) UpdateReason(c *gin.Context, id int64, reason string) error {
	ret := log.db.Table(log.table).Where("id = ?", id).Update("reason", reason)
	return ret.Error
}
