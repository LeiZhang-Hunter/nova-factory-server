package systemservice

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type ISysPermissionService interface {
	SelectPermissionList(c *gin.Context, permission *systemmodels.SysPermissionDQL) (list []*systemmodels.SysPermissionVo)
	SelectPermissionById(ctx context.Context, permissionId int64) (Permission *systemmodels.SysPermissionVo)
	InsertPermission(ctx context.Context, permission *systemmodels.SysPermissionAdd)
	UpdatePermission(ctx context.Context, permission *systemmodels.SysPermissionEdit)
	DeletePermissionById(ctx context.Context, permissionId int64)
	HasChildByPermissionId(ctx context.Context, permissionId int64) bool
	//CheckPermissionExistRole(ctx context.Context, permissionId int64) bool
	SelectPermissionListByRoleIds(ctx context.Context, roleIds []int64) (list []*systemmodels.SysPermissionVo)
}
