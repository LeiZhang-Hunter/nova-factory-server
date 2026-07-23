package systemservice

import (
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type IDeptService interface {
	SelectDeptList(c *gin.Context, dept *modelquery.SysDeptDQL) (list []*modelresponse.SysDeptVo)
	SelectDeptById(c *gin.Context, deptId int64) (dept *modelresponse.SysDeptVo)
	InsertDept(c *gin.Context, dept *modelrequest.SysDeptDML)
	UpdateDept(c *gin.Context, dept *modelrequest.SysDeptDML)
	DeleteDeptById(c *gin.Context, dept int64)
	CheckDeptNameUnique(c *gin.Context, id, parentId int64, deptName string) bool
	HasChildByDeptId(c *gin.Context, deptId int64) bool
	CheckDeptExistUser(c *gin.Context, deptId int64) bool
}
