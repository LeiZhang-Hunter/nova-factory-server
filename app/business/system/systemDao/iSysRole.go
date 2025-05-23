package systemDao

import (
	"context"
	"nova-factory-server/app/business/system/systemModels"
)

type IRoleDao interface {
	SelectRoleList(ctx context.Context, role *systemModels.SysRoleDQL) (roleList []*systemModels.SysRoleVo, total int64)
	SelectRoleAll(ctx context.Context, role *systemModels.SysRoleDQL) (list []*systemModels.SysRoleVo)
	SelectRoleById(ctx context.Context, roleId int64) (role *systemModels.SysRoleVo)
	SelectBasicRolesByUserId(ctx context.Context, userId int64) (roles []*systemModels.SysRole)
	SelectRoleListByUserId(ctx context.Context, userId int64) (list []int64)
	InsertRole(ctx context.Context, sysRole *systemModels.SysRoleDML)
	UpdateRole(ctx context.Context, sysRole *systemModels.SysRoleDML)
	DeleteRoleByIds(ctx context.Context, ids []int64)
	CheckRoleNameUnique(ctx context.Context, roleName string) int64
	SelectAllocatedList(ctx context.Context, user *systemModels.SysRoleAndUserDQL) (list []*systemModels.SysUserVo, total int64)
	SelectUnallocatedList(ctx context.Context, user *systemModels.SysRoleAndUserDQL) (list []*systemModels.SysUserVo, total int64)
	SelectRoleIdAndNameAll(ctx context.Context) (list []*systemModels.SysRoleIdAndName)
	SelectRoleIdAndName(ctx context.Context, userId int64, roleIds []int64) (list []*systemModels.SysRoleIdAndName)
}
