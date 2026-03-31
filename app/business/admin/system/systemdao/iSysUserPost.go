package systemdao

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemmodels"
)

type IUserPostDao interface {
	BatchUserPost(ctx context.Context, users []*systemmodels.SysUserPost)
	DeleteUserPostByUserId(ctx context.Context, userId int64)
	DeleteUserPost(ctx context.Context, ids []int64)
}
