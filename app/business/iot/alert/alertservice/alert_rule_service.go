package alertservice

import (
	alertModels2 "nova-factory-server/app/business/iot/alert/alertmodels"

	"github.com/gin-gonic/gin"
)

type AlertRuleService interface {
	Create(c *gin.Context, data *alertModels2.SetSysAlert) (*alertModels2.SysAlert, error)
	Update(c *gin.Context, data *alertModels2.SetSysAlert) (*alertModels2.SysAlert, error)
	List(c *gin.Context, req *alertModels2.SysAlertListReq) (*alertModels2.SysAlertList, error)
	Remove(c *gin.Context, ids []string) error
	Change(c *gin.Context, data *alertModels2.ChangeSysAlert) error
	FindOpen(c *gin.Context, gatewayId int64) (*alertModels2.SysAlert, error)
	GetReasonByGatewayId(c *gin.Context, gatewayId int64) (*alertModels2.SysAlertAiReason, error)
}
