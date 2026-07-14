package monitorServiceImpl

import (
	"context"
	"fmt"
	"time"

	"nova-factory-server/app/business/admin/monitor/monitordao"
	"nova-factory-server/app/business/admin/monitor/monitormodels"
	"nova-factory-server/app/business/admin/monitor/monitorservice"
	"nova-factory-server/app/utils/store/requestlog"
)

type RequestLogService struct {
	dao monitordao.IRequestLog
}

func NewRequestLogService(dao monitordao.IRequestLog) monitorservice.IRequestLogService {
	svc := &RequestLogService{dao: dao}
	requestlog.RegisterStore(requestLogStore{svc: svc})
	return svc
}

type requestLogStore struct {
	svc *RequestLogService
}

func (s requestLogStore) Create(ctx context.Context, entry *requestlog.Entry) (int64, error) {
	return s.svc.CreateEntry(ctx, entry)
}

func (s requestLogStore) UpdateStatus(ctx context.Context, id int64, status int32, errorMessage string) error {
	return s.svc.UpdateStatus(ctx, id, status, errorMessage)
}

func (s *RequestLogService) CreateEntry(ctx context.Context, entry *requestlog.Entry) (int64, error) {
	data := &monitormodels.RequestLog{
		LogType:       entry.LogType,
		SourceModule:  entry.SourceModule,
		Status:        entry.Status,
		RequestMethod: entry.RequestMethod,
		RequestPath:   entry.RequestPath,
		QueryString:   entry.QueryString,
		HeadersJSON:   entry.HeadersJSON,
		BodyText:      entry.BodyText,
		ClientIP:      entry.ClientIP,
		UserAgent:     entry.UserAgent,
		ErrorMessage:  entry.ErrorMessage,
	}
	if err := s.dao.Create(ctx, data); err != nil {
		return 0, err
	}
	return data.ID, nil
}

func (s *RequestLogService) Create(ctx context.Context, data *monitormodels.RequestLog) (int64, error) {
	if err := s.dao.Create(ctx, data); err != nil {
		return 0, err
	}
	return data.ID, nil
}

func (s *RequestLogService) UpdateStatus(ctx context.Context, id int64, status int32, errorMessage string) error {
	return s.dao.UpdateStatus(ctx, id, status, errorMessage)
}

func (s *RequestLogService) List(ctx context.Context, query *monitormodels.RequestLogQuery) (*monitormodels.RequestLogListData, error) {
	return s.dao.List(ctx, query)
}

func (s *RequestLogService) Detail(ctx context.Context, id int64) (*monitormodels.RequestLog, error) {
	return s.dao.Detail(ctx, id)
}

func (s *RequestLogService) Clean(ctx context.Context, beforeTime string) error {
	if beforeTime == "" {
		beforeTime = time.Now().AddDate(0, 0, -30).Format("2006-01-02 15:04:05")
	}
	if err := s.dao.Clean(ctx, beforeTime); err != nil {
		return fmt.Errorf("clean request log failed: %w", err)
	}
	return nil
}

func (s *RequestLogService) ScheduleClean(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			bg, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			_ = s.Clean(bg, time.Now().AddDate(0, 0, -30).Format("2006-01-02 15:04:05"))
			cancel()
		}
	}
}
