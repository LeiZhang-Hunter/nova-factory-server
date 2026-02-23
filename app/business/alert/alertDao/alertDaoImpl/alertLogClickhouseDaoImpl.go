package alertDaoImpl

import (
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/datasource/clickhouse"

	"github.com/gin-gonic/gin"
)

type AlertLogClickhouseDaoImpl struct {
	clickhouse *clickhouse.ClickHouse
	agentDao   daemonizeDao.IotAgentDao
	table      string
}

func NewAlertLogClickhouseDaoImpl(clickhouse *clickhouse.ClickHouse, agentDao daemonizeDao.IotAgentDao) alertDao.AlertLogClickhouseDao {
	return &AlertLogClickhouseDaoImpl{
		clickhouse: clickhouse,
		table:      "nova_alert_log",
		agentDao:   agentDao,
	}
}

// Export 导出告警数据
func (log *AlertLogClickhouseDaoImpl) Export(c *gin.Context, data []*alertModels.NovaAlertLog) error {
	ret := log.clickhouse.DB().Table(log.table).Create(data)
	if ret.Error != nil {
		return ret.Error
	}

	return nil
}

func (log *AlertLogClickhouseDaoImpl) List(c *gin.Context, req *alertModels.SysAlertLogListReq) (*alertModels.NovaAlertLogList, error) {
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
	var dto []*alertModels.NovaAlertLog

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &alertModels.NovaAlertLogList{
			Rows:  make([]*alertModels.NovaAlertLog, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("start_time_unix desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &alertModels.NovaAlertLogList{
			Rows:  make([]*alertModels.NovaAlertLog, 0),
			Total: 0,
		}, ret.Error
	}

	return &alertModels.NovaAlertLogList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (log *AlertLogClickhouseDaoImpl) GetByObjectId(c *gin.Context, objectId uint64) (*alertModels.NovaAlertLog, error) {
	var dto alertModels.NovaAlertLog
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
