package metricDaoIMpl

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorModel"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/datasource/clickhouse"
	"nova-factory-server/app/utils/time"
	systime "time"
)

type IControlLogDaoImpl struct {
	tableName  string
	clickhouse *clickhouse.ClickHouse
}

func NewIControlLogDaoImpl(clickhouse *clickhouse.ClickHouse) metricDao.IControlLogDao {
	return &IControlLogDaoImpl{
		clickhouse: clickhouse,
		tableName:  "nova_control_log",
	}
}

func (i *IControlLogDaoImpl) Export(ctx context.Context, data []*metricModels.NovaControlLog) error {
	if len(data) == 0 {
		return nil
	}
	ret := i.clickhouse.DB().Table(i.tableName).Create(data)
	if ret.Error != nil {
		zap.L().Error("create device metric data error:", zap.Error(ret.Error))
		return ret.Error
	}
	return ret.Error
}

// List 控制列表
func (i *IControlLogDaoImpl) List(c *gin.Context, req *deviceMonitorModel.ControlLogListReq) (*metricModels.NovaControlLogList, error) {
	if req == nil {
		req = &deviceMonitorModel.ControlLogListReq{}
	}
	db := i.clickhouse.DB().Table(i.tableName)

	if req != nil && req.DeviceID != 0 {
		db = db.Where("device_id = ?", req.DeviceID)
	}
	if req != nil && req.DataId != 0 {
		db = db.Where("data_id = ?", req.DataId)
	}
	if req != nil && req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}

	var startTime string
	if req.Start > 0 {
		startTime = time.GetStartTime(req.Start, 200)
	} else {
		start := systime.Now().UnixMilli() - 86400*1000
		startTime = time.GetStartTime(uint64(start), 200)
	}
	endTime := time.GetEndTimeUseNow(req.End, true)
	if startTime != "" && endTime != "" {
		db = db.Where("time_unix > ?", startTime)
		db = db.Where("time_unix <= ?", endTime)
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
	var dto []*metricModels.NovaControlLog

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &metricModels.NovaControlLogList{
			Rows:  make([]*metricModels.NovaControlLog, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("start_time_unix desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &metricModels.NovaControlLogList{
			Rows:  make([]*metricModels.NovaControlLog, 0),
			Total: 0,
		}, ret.Error
	}

	return &metricModels.NovaControlLogList{
		Rows:  dto,
		Total: total,
	}, nil
}
