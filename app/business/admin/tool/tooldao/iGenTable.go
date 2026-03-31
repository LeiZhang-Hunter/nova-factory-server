package tooldao

import (
	"context"
	"nova-factory-server/app/business/admin/tool/toolmodels"
)

type IGenTable interface {
	SelectGenTableList(ctx context.Context, table *toolmodels.GenTableDQL) (list []*toolmodels.GenTableVo, total int64)
	SelectDbTableList(ctx context.Context, table *toolmodels.GenTableDQL) (list []*toolmodels.DBTableVo, total int64)
	SelectDbTableListByNames(ctx context.Context, tableNames []string) (list []*toolmodels.DBTableVo)
	SelectGenTableById(ctx context.Context, id int64) (table *toolmodels.GenTableVo)
	SelectGenTableAll(ctx context.Context) (list []*toolmodels.GenTableVo)
	BatchInsertGenTable(ctx context.Context, table []*toolmodels.GenTableDML)
	InsertGenTable(ctx context.Context, table *toolmodels.GenTableDML)
	UpdateGenTable(ctx context.Context, table *toolmodels.GenTableDML)
	DeleteGenTableByIds(ctx context.Context, ids []int64)
}
