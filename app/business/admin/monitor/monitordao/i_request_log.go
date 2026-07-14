package monitordao

import (
	"context"
	"nova-factory-server/app/business/admin/monitor/monitormodels"
)

type IRequestLog interface {
	Create(ctx context.Context, data *monitormodels.RequestLog) error
	UpdateStatus(ctx context.Context, id int64, status int32, errorMessage string) error
	List(ctx context.Context, query *monitormodels.RequestLogQuery) (*monitormodels.RequestLogListData, error)
	Detail(ctx context.Context, id int64) (*monitormodels.RequestLog, error)
	Clean(ctx context.Context, beforeTime string) error
}
