package systemdao

import (
	"context"
	modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"
)

type IRoleDao interface {
	SelectRoleList(ctx context.Context, role *modelquery.SysRoleDQL) (roleList []*modelresponse.SysRoleVo, total int64)
	SelectRoleAll(ctx context.Context, role *modelquery.SysRoleDQL) (list []*modelresponse.SysRoleVo)
	SelectRoleById(ctx context.Context, roleId int64) (role *modelresponse.SysRoleVo)
	SelectBasicRolesByUserId(ctx context.Context, userId int64) (roles []*modelentity.SysRole)
	SelectRoleListByUserId(ctx context.Context, userId int64) (list []int64)
	InsertRole(ctx context.Context, sysRole *modelrequest.SysRoleDML)
	UpdateRole(ctx context.Context, sysRole *modelrequest.SysRoleDML)
	DeleteRoleByIds(ctx context.Context, ids []int64)
	CheckRoleNameUnique(ctx context.Context, roleName string) int64
	SelectAllocatedList(ctx context.Context, user *modelquery.SysRoleAndUserDQL) (list []*modelresponse.SysUserVo, total int64)
	SelectUnallocatedList(ctx context.Context, user *modelquery.SysRoleAndUserDQL) (list []*modelresponse.SysUserVo, total int64)
	SelectRoleIdAndNameAll(ctx context.Context) (list []*modelresponse.SysRoleIdAndName)
	SelectRoleIdAndName(ctx context.Context, userId int64, roleIds []int64) (list []*modelresponse.SysRoleIdAndName)
}
