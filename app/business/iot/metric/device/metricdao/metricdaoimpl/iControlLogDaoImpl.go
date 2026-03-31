package metricdaoimpl

import (
	"context"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitormodel"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"
	"nova-factory-server/app/datasource/clickhouse"
	"nova-factory-server/app/utils/time"
	systime "time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IControlLogDaoImpl struct {
	tableName  string
	clickhouse *clickhouse.ClickHouse
}

func NewIControlLogDaoImpl(clickhouse *clickhouse.ClickHouse) metricdao.IControlLogDao {
	return &IControlLogDaoImpl{
		clickhouse: clickhouse,
		tableName:  "nova_control_log",
	}
}

func (i *IControlLogDaoImpl) Export(ctx context.Context, data []*metricmodels.NovaControlLog) error {
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
func (i *IControlLogDaoImpl) List(c *gin.Context, req *devicemonitormodel.ControlLogListReq) (*metricmodels.NovaControlLogList, error) {
	if req == nil {
		req = &devicemonitormodel.ControlLogListReq{}
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
	var dto []*metricmodels.NovaControlLog

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &metricmodels.NovaControlLogList{
			Rows:  make([]*metricmodels.NovaControlLog, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("start_time_unix desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &metricmodels.NovaControlLogList{
			Rows:  make([]*metricmodels.NovaControlLog, 0),
			Total: 0,
		}, ret.Error
	}

	return &metricmodels.NovaControlLogList{
		Rows:  dto,
		Total: total,
	}, nil
}
