package systemdao

import (
	"context"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"
)

type IPostDao interface {
	SelectPostAll(ctx context.Context) (sysPost []*modelresponse.SysPostVo)
	SelectPostListByUserId(ctx context.Context, userId int64) (list []int64)
	SelectPostList(ctx context.Context, post *modelquery.SysPostDQL) (list []*modelresponse.SysPostVo, total int64)
	SelectPostListAll(ctx context.Context, post *modelquery.SysPostDQL) (list []*modelresponse.SysPostVo)
	SelectPostById(ctx context.Context, postId int64) (dictData *modelresponse.SysPostVo)
	InsertPost(ctx context.Context, post *modelrequest.SysPostDML)
	UpdatePost(ctx context.Context, post *modelrequest.SysPostDML)
	DeletePostByIds(ctx context.Context, dictCodes []int64)
	SelectPostNameListByUserId(ctx context.Context, userId int64) (list []string)
}
