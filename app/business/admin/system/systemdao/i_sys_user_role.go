package systemdao

import (
	"context"
	modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"
)

type IUserRoleDao interface {
	DeleteUserRole(ctx context.Context, ids []int64)
	BatchUserRole(ctx context.Context, users []*modelentity.SysUserRole)
	DeleteUserRoleByUserId(ctx context.Context, userId int64)
	CountUserRoleByRoleId(ctx context.Context, ids []int64) int
	DeleteUserRoleInfo(ctx context.Context, userRole *modelentity.SysUserRole)
	DeleteUserRoleInfos(ctx context.Context, roleId int64, userIds []int64)
}
