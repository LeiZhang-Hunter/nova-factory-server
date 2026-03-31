package alertservice

import (
	"nova-factory-server/app/business/iot/alert/alertmodels"

	"github.com/gin-gonic/gin"
)

type AlertActionService interface {
	Set(c *gin.Context, data *alertmodels.SetAlertAction) (*alertmodels.AlertAction, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *alertmodels.SysAlertActionListReq) (*alertmodels.SysAlertActionList, error)
}
