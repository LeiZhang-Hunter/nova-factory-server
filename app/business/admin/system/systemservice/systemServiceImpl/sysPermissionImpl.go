package systemServiceImpl

import (
	"context"
	"nova-factory-server/app/business/admin/system/systemdao"
	"nova-factory-server/app/business/admin/system/systemmodels"
	"nova-factory-server/app/business/admin/system/systemservice"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
)

type PermissionService struct {
	pd systemdao.IPermissionDao
}

func NewPermissionService(pd systemdao.IPermissionDao) systemservice.ISysPermissionService {
	return &PermissionService{pd: pd}
}
func (ps *PermissionService) SelectPermissionList(c *gin.Context, permission *systemmodels.SysPermissionDQL) (list []*systemmodels.SysPermissionVo) {

	list = ps.pd.SelectPermissionList(c, permission)
	return
}

func (ps *PermissionService) SelectPermissionById(ctx context.Context, permissionId int64) (Permission *systemmodels.SysPermissionVo) {
	return ps.pd.SelectPermissionById(ctx, permissionId)
}

func (ps *PermissionService) SelectPermissionListByRoleIds(ctx context.Context, roleIds []int64) (list []*systemmodels.SysPermissionVo) {
	return ps.pd.SelectPermissionListByRoleIds(ctx, roleIds)
}

func (ps *PermissionService) InsertPermission(ctx context.Context, permission *systemmodels.SysPermissionAdd) {
	permission.PermissionId = snowflake.GenID()
	permission.Status = "0"
	ps.pd.InsertPermission(ctx, permission)
}

func (ps *PermissionService) UpdatePermission(ctx context.Context, permission *systemmodels.SysPermissionEdit) {
	ps.pd.UpdatePermission(ctx, permission)
}

func (ps *PermissionService) DeletePermissionById(ctx context.Context, permissionId int64) {
	ps.pd.DeletePermissionById(ctx, permissionId)
}

func (ps *PermissionService) HasChildByPermissionId(ctx context.Context, permissionId int64) bool {
	return ps.pd.HasChildByPermissionId(ctx, permissionId) > 0
}

//func (ps *PermissionService) CheckPermissionExistRole(ctx context.Context, permissionId int64) bool {
//	return ps.rd.CheckPermissionExistRole(ctx, permissionId) > 0
//}
