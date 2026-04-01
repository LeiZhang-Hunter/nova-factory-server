package tooldao

import (
	"context"
	"nova-factory-server/app/business/admin/tool/toolmodels"
)

type IGenTableColumn interface {
	SelectDbTableColumnsByName(ctx context.Context, tableName string) (list []*toolmodels.InformationSchemaColumn)
	SelectGenTableColumnListByTableId(ctx context.Context, tableId int64) (list []*toolmodels.GenTableColumnVo)
	BatchInsertGenTableColumn(ctx context.Context, genTables []*toolmodels.GenTableColumnDML)
	UpdateGenTableColumn(ctx context.Context, column *toolmodels.GenTableColumnDML)
	DeleteGenTableColumnByIds(ctx context.Context, ids []int64)
}
