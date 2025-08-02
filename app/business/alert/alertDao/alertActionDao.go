package alertDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertModels"
)

type AlertActionDao interface {
	Set(c *gin.Context, data *alertModels.SetAlertAction) (*alertModels.AlertAction, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *alertModels.SysAlertActionListReq) (*alertModels.SysAlertActionList, error)
}
