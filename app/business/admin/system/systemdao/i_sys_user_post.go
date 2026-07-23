package systemdao

import (
	"context"
	modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"
)

type IUserPostDao interface {
	BatchUserPost(ctx context.Context, users []*modelentity.SysUserPost)
	DeleteUserPostByUserId(ctx context.Context, userId int64)
	DeleteUserPost(ctx context.Context, ids []int64)
}
