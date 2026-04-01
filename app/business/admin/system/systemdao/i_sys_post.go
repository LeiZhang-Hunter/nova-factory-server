package systemdao

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemmodels"
)

type IPostDao interface {
	SelectPostAll(ctx context.Context) (sysPost []*systemmodels.SysPostVo)
	SelectPostListByUserId(ctx context.Context, userId int64) (list []int64)
	SelectPostList(ctx context.Context, post *systemmodels.SysPostDQL) (list []*systemmodels.SysPostVo, total int64)
	SelectPostListAll(ctx context.Context, post *systemmodels.SysPostDQL) (list []*systemmodels.SysPostVo)
	SelectPostById(ctx context.Context, postId int64) (dictData *systemmodels.SysPostVo)
	InsertPost(ctx context.Context, post *systemmodels.SysPostVo)
	UpdatePost(ctx context.Context, post *systemmodels.SysPostVo)
	DeletePostByIds(ctx context.Context, dictCodes []int64)
	SelectPostNameListByUserId(ctx context.Context, userId int64) (list []string)
}
