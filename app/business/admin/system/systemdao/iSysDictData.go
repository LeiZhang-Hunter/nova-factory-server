package systemdao

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemmodels"
)

type IDictDataDao interface {
	SelectDictDataByType(ctx context.Context, dictType string) (SysDictDataList []*systemmodels.SysDictDataVo)
	SelectDictDataList(ctx context.Context, dictData *systemmodels.SysDictDataDQL) (list []*systemmodels.SysDictDataVo, total int64)
	SelectDictDataById(ctx context.Context, dictCode int64) (dictData *systemmodels.SysDictDataVo)
	InsertDictData(ctx context.Context, dictData *systemmodels.SysDictDataVo)
	UpdateDictData(ctx context.Context, dictData *systemmodels.SysDictDataVo)
	SelectDictTypesByDictCodes(ctx context.Context, dictCodes []int64) []string
	DeleteDictDataByIds(ctx context.Context, dictCodes []int64)
	CountDictDataByTypes(ctx context.Context, dictType []string) int
}
