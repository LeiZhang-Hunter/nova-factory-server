package alertDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
)

type AlertLogDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewAlertLogDaoImpl(db *gorm.DB) alertDao.AlertLogDao {
	return &AlertLogDaoImpl{
		db:    db,
		table: "sys_alert_log",
	}
}

func (log *AlertLogDaoImpl) Export(c *gin.Context, data []*alertModels.SysAlertLog) error {
	ret := log.db.Table(log.table).Create(data)
	return ret.Error
}

func (log *AlertLogDaoImpl) List(c *gin.Context, req *alertModels.SysAlertLogListReq) (*alertModels.SysAlertLogList, error) {
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
	var dto []*alertModels.SysAlertLog

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &alertModels.SysAlertLogList{
			Rows:  make([]*alertModels.SysAlertLog, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &alertModels.SysAlertLogList{
			Rows:  make([]*alertModels.SysAlertLog, 0),
			Total: 0,
		}, ret.Error
	}

	return &alertModels.SysAlertLogList{
		Rows:  dto,
		Total: total,
	}, nil
}
