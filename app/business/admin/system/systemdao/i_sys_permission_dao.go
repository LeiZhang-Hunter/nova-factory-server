package systemdao

import (
	"context"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"
)

type IPermissionDao interface {
	SelectPermissionByUserId(ctx context.Context, userId int64) []string
	SelectPermissionList(ctx context.Context, permission *modelquery.SysPermissionDQL) (list []*modelresponse.SysPermissionVo)
	SelectPermissionById(ctx context.Context, permissionId int64) *modelresponse.SysPermissionVo
	SelectPermissionListByParentId(ctx context.Context, parentId int64) (list []*modelresponse.SysPermissionVo)
	SelectPermissionListByRoleIds(ctx context.Context, roleIds []int64) (list []*modelresponse.SysPermissionVo)
	InsertPermission(ctx context.Context, permission *modelrequest.SysPermissionAdd)
	UpdatePermission(ctx context.Context, permission *modelrequest.SysPermissionEdit)
	DeletePermissionById(ctx context.Context, permissionId int64)
	HasChildByPermissionId(ctx context.Context, permissionId int64) int
	SelectPermissionAll(ctx context.Context) []string
	SelectPermissionListSelectBoxByPerm(ctx context.Context, perm []string) (list []*modelresponse.SelectPermission)
}
