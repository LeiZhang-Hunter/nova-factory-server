package alertdao

import (
	"nova-factory-server/app/business/iot/alert/alertmodels"

	"github.com/gin-gonic/gin"
)

type AlertAiReasonDao interface {
	Set(c *gin.Context, data *alertmodels.SetAlertAiReason) (*alertmodels.SysAlertAiReason, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *alertmodels.SysAlertAiReasonReq) (*alertmodels.SysAlertReasonList, error)
	GetById(c *gin.Context, id int64) (*alertmodels.SysAlertAiReason, error)
}
