package systemservice

import (
	"nova-factory-server/app/baize"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type ISelectBoxService interface {
	SelectPermissionBox(c *gin.Context) (list []*modelresponse.SelectPermission)
	SelectDeptBox(c *gin.Context, be *baize.BaseEntityDQL) (list []*modelresponse.SelectDept)
}
