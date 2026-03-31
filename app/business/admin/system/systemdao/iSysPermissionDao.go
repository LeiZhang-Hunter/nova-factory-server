package systemdao

import (
	"context"
	systemModels2 "nova-factory-server/app/business/admin/system/systemmodels"
)

type IPermissionDao interface {
	SelectPermissionByUserId(ctx context.Context, userId int64) []string
	SelectPermissionList(ctx context.Context, permission *systemModels2.SysPermissionDQL) (list []*systemModels2.SysPermissionVo)
	SelectPermissionById(ctx context.Context, permissionId int64) *systemModels2.SysPermissionVo
	SelectPermissionListByParentId(ctx context.Context, parentId int64) (list []*systemModels2.SysPermissionVo)
	SelectPermissionListByRoleIds(ctx context.Context, roleIds []int64) (list []*systemModels2.SysPermissionVo)
	InsertPermission(ctx context.Context, permission *systemModels2.SysPermissionAdd)
	UpdatePermission(ctx context.Context, permission *systemModels2.SysPermissionEdit)
	DeletePermissionById(ctx context.Context, permissionId int64)
	HasChildByPermissionId(ctx context.Context, permissionId int64) int
	SelectPermissionAll(ctx context.Context) []string
	SelectPermissionListSelectBoxByPerm(ctx context.Context, perm []string) (list []*systemModels2.SelectPermission)
}
