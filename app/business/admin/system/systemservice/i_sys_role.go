package systemservice

import (
	modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type IRoleService interface {
	SelectRoleList(c *gin.Context, role *modelquery.SysRoleDQL) (list []*modelresponse.SysRoleVo, total int64)
	RoleExport(c *gin.Context, role *modelquery.SysRoleDQL) (data []byte)
	SelectRoleById(c *gin.Context, roseId int64) (role *modelresponse.SysRoleVo)
	InsertRole(c *gin.Context, sysRole *modelrequest.SysRoleDML)
	UpdateRole(c *gin.Context, sysRole *modelrequest.SysRoleDML)
	UpdateRoleStatus(c *gin.Context, sysRole *modelrequest.SysRoleDML)
	DeleteRoleByIds(c *gin.Context, ids []int64)
	CountUserRoleByRoleId(c *gin.Context, ids []int64) bool
	SelectBasicRolesByUserId(c *gin.Context, userId int64) (roles []*modelentity.SysRole)
	RolePermissionByRoles(c *gin.Context, roles []*modelentity.SysRole) (loginRoles []int64)

	CheckRoleNameUnique(c *gin.Context, id int64, roleName string) bool
	SelectAllocatedList(c *gin.Context, user *modelquery.SysRoleAndUserDQL) (list []*modelresponse.SysUserVo, total int64)
	SelectUnallocatedList(c *gin.Context, user *modelquery.SysRoleAndUserDQL) (list []*modelresponse.SysUserVo, total int64)
	InsertAuthUsers(c *gin.Context, roleId int64, userIds []int64)
	DeleteAuthUsers(c *gin.Context, roleId int64, userIds []int64)
	DeleteAuthUserRole(c *gin.Context, user *modelentity.SysUserRole)
}
