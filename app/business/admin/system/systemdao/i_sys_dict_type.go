package systemdao

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemmodels"
)

type IDictTypeDao interface {
	SelectDictTypeList(ctx context.Context, dictType *systemmodels.SysDictTypeDQL) (list []*systemmodels.SysDictTypeVo, total int64)
	SelectDictTypeAll(ctx context.Context, dictType *systemmodels.SysDictTypeDQL) (list []*systemmodels.SysDictTypeVo)
	SelectDictTypeById(ctx context.Context, dictId int64) (dictType *systemmodels.SysDictTypeVo)
	SelectDictTypeByIds(ctx context.Context, dictId []int64) (dictTypes []string)
	InsertDictType(ctx context.Context, dictType *systemmodels.SysDictTypeVo)
	UpdateDictType(ctx context.Context, dictType *systemmodels.SysDictTypeVo)
	DeleteDictTypeByIds(ctx context.Context, dictIds []int64)
	CheckDictTypeUnique(ctx context.Context, dictType string) int64
}
