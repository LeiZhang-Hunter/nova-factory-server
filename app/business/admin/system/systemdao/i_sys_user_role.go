package systemdao

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemmodels"
)

type IUserRoleDao interface {
	DeleteUserRole(ctx context.Context, ids []int64)
	BatchUserRole(ctx context.Context, users []*systemmodels.SysUserRole)
	DeleteUserRoleByUserId(ctx context.Context, userId int64)
	CountUserRoleByRoleId(ctx context.Context, ids []int64) int
	DeleteUserRoleInfo(ctx context.Context, userRole *systemmodels.SysUserRole)
	DeleteUserRoleInfos(ctx context.Context, roleId int64, userIds []int64)
}
