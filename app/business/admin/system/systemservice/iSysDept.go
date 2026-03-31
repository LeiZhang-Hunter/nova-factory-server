package systemservice

import (
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type IDeptService interface {
	SelectDeptList(c *gin.Context, dept *systemmodels.SysDeptDQL) (list []*systemmodels.SysDeptVo)
	SelectDeptById(c *gin.Context, deptId int64) (dept *systemmodels.SysDeptVo)
	InsertDept(c *gin.Context, dept *systemmodels.SysDeptVo)
	UpdateDept(c *gin.Context, dept *systemmodels.SysDeptVo)
	DeleteDeptById(c *gin.Context, dept int64)
	CheckDeptNameUnique(c *gin.Context, id, parentId int64, deptName string) bool
	HasChildByDeptId(c *gin.Context, deptId int64) bool
	CheckDeptExistUser(c *gin.Context, deptId int64) bool
}
