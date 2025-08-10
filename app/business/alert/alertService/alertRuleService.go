package alertService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertModels"
)

type AlertRuleService interface {
	Create(c *gin.Context, data *alertModels.SetSysAlert) (*alertModels.SysAlert, error)
	Update(c *gin.Context, data *alertModels.SetSysAlert) (*alertModels.SysAlert, error)
	List(c *gin.Context, req *alertModels.SysAlertListReq) (*alertModels.SysAlertList, error)
	Remove(c *gin.Context, ids []string) error
	Change(c *gin.Context, data *alertModels.ChangeSysAlert) error
	FindOpen(c *gin.Context) (*alertModels.SysAlert, error)
}
