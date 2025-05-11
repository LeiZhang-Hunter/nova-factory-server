package systemService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/system/systemModels"
)

type ISelectBoxService interface {
	SelectPermissionBox(c *gin.Context) (list []*systemModels.SelectPermission)
	SelectDeptBox(c *gin.Context, be *baize.BaseEntityDQL) (list []*systemModels.SelectDept)
}
