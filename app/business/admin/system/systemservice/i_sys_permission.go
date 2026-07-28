package systemservice

import (
	"context"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type ISysPermissionService interface {
	SelectPermissionList(c *gin.Context, permission *modelquery.SysPermissionDQL) (list []*modelresponse.SysPermissionVo)
	SelectPermissionById(ctx context.Context, permissionId int64) (Permission *modelresponse.SysPermissionVo)
	InsertPermission(ctx context.Context, permission *modelrequest.SysPermissionAdd)
	UpdatePermission(ctx context.Context, permission *modelrequest.SysPermissionEdit)
	DeletePermissionById(ctx context.Context, permissionId int64)
	HasChildByPermissionId(ctx context.Context, permissionId int64) bool
	//CheckPermissionExistRole(ctx context.Context, permissionId int64) bool
	SelectPermissionListByRoleIds(ctx context.Context, roleIds []int64) (list []*modelresponse.SysPermissionVo)
}
