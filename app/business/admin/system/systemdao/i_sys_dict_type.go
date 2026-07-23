package systemdao

import (
	"context"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"
)

type IDictTypeDao interface {
	SelectDictTypeList(ctx context.Context, dictType *modelquery.SysDictTypeDQL) (list []*modelresponse.SysDictTypeVo, total int64)
	SelectDictTypeAll(ctx context.Context, dictType *modelquery.SysDictTypeDQL) (list []*modelresponse.SysDictTypeVo)
	SelectDictTypeById(ctx context.Context, dictId int64) (dictType *modelresponse.SysDictTypeVo)
	SelectDictTypeByIds(ctx context.Context, dictId []int64) (dictTypes []string)
	InsertDictType(ctx context.Context, dictType *modelrequest.SysDictTypeDML)
	UpdateDictType(ctx context.Context, dictType *modelrequest.SysDictTypeDML)
	DeleteDictTypeByIds(ctx context.Context, dictIds []int64)
	CheckDictTypeUnique(ctx context.Context, dictType string) int64
}
