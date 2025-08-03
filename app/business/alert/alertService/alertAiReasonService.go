package alertService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertModels"
)

type AlertAiReasonService interface {
	Set(c *gin.Context, data *alertModels.SetAlertAiReason) (*alertModels.SysAlertAiReason, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *alertModels.SysAlertAiReasonReq) (*alertModels.SysAlertReasonList, error)
}
