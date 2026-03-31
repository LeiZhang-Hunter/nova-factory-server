package alertDao

import (
	"nova-factory-server/app/business/iot/alert/alertModels"

	"github.com/gin-gonic/gin"
)

type AlertActionDao interface {
	Set(c *gin.Context, data *alertModels.SetAlertAction) (*alertModels.AlertAction, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *alertModels.SysAlertActionListReq) (*alertModels.SysAlertActionList, error)
	GetById(c *gin.Context, id int64) (*alertModels.SetAlertAction, error)
}
