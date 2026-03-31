package alertDao

import (
	"nova-factory-server/app/business/iot/alert/alertModels"

	"github.com/gin-gonic/gin"
)

type AlertAiReasonDao interface {
	Set(c *gin.Context, data *alertModels.SetAlertAiReason) (*alertModels.SysAlertAiReason, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *alertModels.SysAlertAiReasonReq) (*alertModels.SysAlertReasonList, error)
	GetById(c *gin.Context, id int64) (*alertModels.SysAlertAiReason, error)
}
