package monitordao

import (
	"context"
	"nova-factory-server/app/business/admin/monitor/monitormodels"
)

type ILogininforDao interface {
	InserLogininfor(ctx context.Context, logininfor *monitormodels.Logininfor)
	SelectLogininforList(ctx context.Context, logininfor *monitormodels.LogininforDQL) (list []*monitormodels.Logininfor, total int64)
	SelectLogininforListAll(ctx context.Context, logininfor *monitormodels.LogininforDQL) (list []*monitormodels.Logininfor)
	DeleteLogininforByIds(ctx context.Context, infoIds []int64)
	CleanLogininfor(ctx context.Context)
}
