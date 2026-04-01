package systemservice

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type ISelectBoxService interface {
	SelectPermissionBox(c *gin.Context) (list []*systemmodels.SelectPermission)
	SelectDeptBox(c *gin.Context, be *baize.BaseEntityDQL) (list []*systemmodels.SelectDept)
}
