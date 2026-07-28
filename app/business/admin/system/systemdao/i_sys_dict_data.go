package systemdao

import (
	"context"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"
)

type IDictDataDao interface {
	SelectDictDataByType(ctx context.Context, dictType string) (SysDictDataList []*modelresponse.SysDictDataVo)
	SelectDictDataList(ctx context.Context, dictData *modelquery.SysDictDataDQL) (list []*modelresponse.SysDictDataVo, total int64)
	SelectDictDataById(ctx context.Context, dictCode int64) (dictData *modelresponse.SysDictDataVo)
	InsertDictData(ctx context.Context, dictData *modelrequest.SysDictDataDML)
	UpdateDictData(ctx context.Context, dictData *modelrequest.SysDictDataDML)
	SelectDictTypesByDictCodes(ctx context.Context, dictCodes []int64) []string
	DeleteDictDataByIds(ctx context.Context, dictCodes []int64)
	CountDictDataByTypes(ctx context.Context, dictType []string) int
}
