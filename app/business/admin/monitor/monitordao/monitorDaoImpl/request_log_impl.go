package monitorDaoImpl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/business/admin/monitor/monitordao"
	"nova-factory-server/app/business/admin/monitor/monitormodels"

	"gorm.io/gorm"
)

type RequestLogDao struct {
	db        *gorm.DB
	tableName string
}

func NewRequestLogDao(db *gorm.DB) monitordao.IRequestLog {
	return &RequestLogDao{db: db, tableName: "sys_request_log"}
}

func (d *RequestLogDao) Create(ctx context.Context, data *monitormodels.RequestLog) error {
	if data == nil {
		return errors.New("request log is nil")
	}
	return d.db.WithContext(ctx).Table(d.tableName).Create(data).Error
}

func (d *RequestLogDao) UpdateStatus(ctx context.Context, id int64, status int32, errorMessage string) error {
	if id <= 0 {
		return nil
	}
	updates := map[string]any{
		"status":      status,
		"update_time": time.Now(),
	}
	if strings.TrimSpace(errorMessage) != "" {
		updates["error_message"] = errorMessage
	}
	return d.db.WithContext(ctx).Table(d.tableName).Where("id = ?", id).Updates(updates).Error
}

func (d *RequestLogDao) List(ctx context.Context, query *monitormodels.RequestLogQuery) (*monitormodels.RequestLogListData, error) {
	db := d.db.WithContext(ctx).Table(d.tableName)
	if query.LogType != "" {
		db = db.Where("log_type = ?", query.LogType)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.BeginTime != "" {
		db = db.Where("create_time >= ?", query.BeginTime)
	}
	if query.EndTime != "" {
		db = db.Where("create_time <= ?", query.EndTime)
	}
	if query.PageNum <= 0 {
		query.PageNum = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*monitormodels.RequestLog, 0)
	if err := db.Order("id DESC").Offset((query.PageNum - 1) * query.PageSize).Limit(query.PageSize).Find(&rows).Error; err != nil {
		return nil, err
	}
	return &monitormodels.RequestLogListData{Rows: rows, Total: total}, nil
}

func (d *RequestLogDao) Detail(ctx context.Context, id int64) (*monitormodels.RequestLog, error) {
	var row monitormodels.RequestLog
	if err := d.db.WithContext(ctx).Table(d.tableName).Where("id = ?", id).First(&row).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

func (d *RequestLogDao) Clean(ctx context.Context, beforeTime string) error {
	if strings.TrimSpace(beforeTime) == "" {
		return fmt.Errorf("beforeTime required")
	}
	return d.db.WithContext(ctx).Table(d.tableName).Where("create_time < ?", beforeTime).Delete(&monitormodels.RequestLog{}).Error
}
