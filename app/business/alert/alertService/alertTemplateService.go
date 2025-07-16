package alertService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertModels"
)

type AlertTemplateService interface {
	Create(c *gin.Context, data *alertModels.SetSysAlertSinkTemplate) (*alertModels.SysAlertSinkTemplate, error)
	Update(c *gin.Context, data *alertModels.SetSysAlertSinkTemplate) (*alertModels.SysAlertSinkTemplate, error)
	List(c *gin.Context, req *alertModels.SysAlertSinkTemplateReq) (*alertModels.SysAlertSinkTemplateListData, error)
	Remove(c *gin.Context, ids []string) error
}
