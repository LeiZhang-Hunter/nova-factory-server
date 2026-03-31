package alertdaoimpl

import (
	"nova-factory-server/app/business/iot/alert/alertdao"
	alertModels2 "nova-factory-server/app/business/iot/alert/alertmodels"
	"nova-factory-server/app/business/iot/daemonize/daemonizedao"
	"nova-factory-server/app/datasource/clickhouse"

	"github.com/gin-gonic/gin"
)

type AlertLogClickhouseDaoImpl struct {
	clickhouse *clickhouse.ClickHouse
	agentDao   daemonizedao.IotAgentDao
	table      string
}

func NewAlertLogClickhouseDaoImpl(clickhouse *clickhouse.ClickHouse, agentDao daemonizedao.IotAgentDao) alertdao.AlertLogClickhouseDao {
	return &AlertLogClickhouseDaoImpl{
		clickhouse: clickhouse,
		table:      "nova_alert_log",
		agentDao:   agentDao,
	}
}

// Export 导出告警数据
func (log *AlertLogClickhouseDaoImpl) Export(c *gin.Context, data []*alertModels2.NovaAlertLog) error {
	ret := log.clickhouse.DB().Table(log.table).Create(data)
	if ret.Error != nil {
		return ret.Error
	}

	return nil
}

func (log *AlertLogClickhouseDaoImpl) List(c *gin.Context, req *alertModels2.SysAlertLogListReq) (*alertModels2.NovaAlertLogList, error) {
	db := log.clickhouse.DB().Table(log.table)
	if req != nil && req.AlertID != 0 {
		db = db.Where("alert_id = ?", req.AlertID)
	}
	if req != nil && req.GatewayID != 0 {
		db = db.Where("gateway_id = ?", req.GatewayID)
	}
	if req != nil && req.DeviceId != 0 {
		db = db.Where("device_id = ?", req.DeviceId)
	}
	if req != nil && req.StartTime != "" {
		db = db.Where("start_time_unix >= ?", req.StartTime)
	}
	if req != nil && req.EndTime != "" {
		db = db.Where("start_time_unix < ?", req.EndTime)
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
	var dto []*alertModels2.NovaAlertLog

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &alertModels2.NovaAlertLogList{
			Rows:  make([]*alertModels2.NovaAlertLog, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("start_time_unix desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &alertModels2.NovaAlertLogList{
			Rows:  make([]*alertModels2.NovaAlertLog, 0),
			Total: 0,
		}, ret.Error
	}

	return &alertModels2.NovaAlertLogList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (log *AlertLogClickhouseDaoImpl) GetByObjectId(c *gin.Context, objectId uint64) (*alertModels2.NovaAlertLog, error) {
	var dto alertModels2.NovaAlertLog
	ret := log.clickhouse.DB().Table(log.table).Where("object_id = ?", objectId).First(&dto)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &dto, nil
}

func (log *AlertLogClickhouseDaoImpl) Count(c *gin.Context) (int64, error) {
	var count int64
	ret := log.clickhouse.DB().Table(log.table).Count(&count)
	if ret.Error != nil {
		return 0, ret.Error
	}
	return count, nil
}
