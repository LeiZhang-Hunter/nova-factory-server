package systemDao

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemModels"
)

type IUserPostDao interface {
	BatchUserPost(ctx context.Context, users []*systemModels.SysUserPost)
	DeleteUserPostByUserId(ctx context.Context, userId int64)
	DeleteUserPost(ctx context.Context, ids []int64)
}
