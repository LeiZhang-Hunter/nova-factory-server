package systemService

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/admin/system/systemModels"

	"github.com/gin-gonic/gin"
)

type ISelectBoxService interface {
	SelectPermissionBox(c *gin.Context) (list []*systemModels.SelectPermission)
	SelectDeptBox(c *gin.Context, be *baize.BaseEntityDQL) (list []*systemModels.SelectDept)
}
