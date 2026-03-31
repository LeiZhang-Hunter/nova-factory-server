package systemservice

import (
	systemModels2 "nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type IRoleService interface {
	SelectRoleList(c *gin.Context, role *systemModels2.SysRoleDQL) (list []*systemModels2.SysRoleVo, total int64)
	RoleExport(c *gin.Context, role *systemModels2.SysRoleDQL) (data []byte)
	SelectRoleById(c *gin.Context, roseId int64) (role *systemModels2.SysRoleVo)
	InsertRole(c *gin.Context, sysRole *systemModels2.SysRoleDML)
	UpdateRole(c *gin.Context, sysRole *systemModels2.SysRoleDML)
	UpdateRoleStatus(c *gin.Context, sysRole *systemModels2.SysRoleDML)
	DeleteRoleByIds(c *gin.Context, ids []int64)
	CountUserRoleByRoleId(c *gin.Context, ids []int64) bool
	SelectBasicRolesByUserId(c *gin.Context, userId int64) (roles []*systemModels2.SysRole)
	RolePermissionByRoles(c *gin.Context, roles []*systemModels2.SysRole) (loginRoles []int64)

	CheckRoleNameUnique(c *gin.Context, id int64, roleName string) bool
	SelectAllocatedList(c *gin.Context, user *systemModels2.SysRoleAndUserDQL) (list []*systemModels2.SysUserVo, total int64)
	SelectUnallocatedList(c *gin.Context, user *systemModels2.SysRoleAndUserDQL) (list []*systemModels2.SysUserVo, total int64)
	InsertAuthUsers(c *gin.Context, roleId int64, userIds []int64)
	DeleteAuthUsers(c *gin.Context, roleId int64, userIds []int64)
	DeleteAuthUserRole(c *gin.Context, user *systemModels2.SysUserRole)
}
