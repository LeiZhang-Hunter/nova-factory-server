package systemdao

import (
	"context"
	systemModels2 "nova-factory-server/app/business/admin/system/systemmodels"
)

type IRoleDao interface {
	SelectRoleList(ctx context.Context, role *systemModels2.SysRoleDQL) (roleList []*systemModels2.SysRoleVo, total int64)
	SelectRoleAll(ctx context.Context, role *systemModels2.SysRoleDQL) (list []*systemModels2.SysRoleVo)
	SelectRoleById(ctx context.Context, roleId int64) (role *systemModels2.SysRoleVo)
	SelectBasicRolesByUserId(ctx context.Context, userId int64) (roles []*systemModels2.SysRole)
	SelectRoleListByUserId(ctx context.Context, userId int64) (list []int64)
	InsertRole(ctx context.Context, sysRole *systemModels2.SysRoleDML)
	UpdateRole(ctx context.Context, sysRole *systemModels2.SysRoleDML)
	DeleteRoleByIds(ctx context.Context, ids []int64)
	CheckRoleNameUnique(ctx context.Context, roleName string) int64
	SelectAllocatedList(ctx context.Context, user *systemModels2.SysRoleAndUserDQL) (list []*systemModels2.SysUserVo, total int64)
	SelectUnallocatedList(ctx context.Context, user *systemModels2.SysRoleAndUserDQL) (list []*systemModels2.SysUserVo, total int64)
	SelectRoleIdAndNameAll(ctx context.Context) (list []*systemModels2.SysRoleIdAndName)
	SelectRoleIdAndName(ctx context.Context, userId int64, roleIds []int64) (list []*systemModels2.SysRoleIdAndName)
}
