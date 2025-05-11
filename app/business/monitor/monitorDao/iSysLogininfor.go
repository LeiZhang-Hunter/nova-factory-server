package monitorDao

import (
	"context"
	"nova-factory-server/app/business/monitor/monitorModels"
)

type ILogininforDao interface {
	InserLogininfor(ctx context.Context, logininfor *monitorModels.Logininfor)
	SelectLogininforList(ctx context.Context, logininfor *monitorModels.LogininforDQL) (list []*monitorModels.Logininfor, total int64)
	SelectLogininforListAll(ctx context.Context, logininfor *monitorModels.LogininforDQL) (list []*monitorModels.Logininfor)
	DeleteLogininforByIds(ctx context.Context, infoIds []int64)
	CleanLogininfor(ctx context.Context)
}
