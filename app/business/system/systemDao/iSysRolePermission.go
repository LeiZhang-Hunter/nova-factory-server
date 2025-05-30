package systemDao

import (
	"context"
	"nova-factory-server/app/business/system/systemModels"
)

type IRolePermissionDao interface {
	SelectPermissionIdsByRoleId(ctx context.Context, roleId int64) []int64
	BatchRolePermission(ctx context.Context, list []*systemModels.SysRolePermission)
	DeleteRolePermission(ctx context.Context, ids []int64)
	DeleteRolePermissionByRoleId(ctx context.Context, roleId int64)
	CheckPermissionExistRole(ctx context.Context, permissionId int64) int
}
