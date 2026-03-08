package alertService

import (
	"nova-factory-server/app/business/iot/alert/alertModels"

	"github.com/gin-gonic/gin"
)

type AlertTemplateService interface {
	Create(c *gin.Context, data *alertModels.SetSysAlertSinkTemplate) (*alertModels.SysAlertSinkTemplate, error)
	Update(c *gin.Context, data *alertModels.SetSysAlertSinkTemplate) (*alertModels.SysAlertSinkTemplate, error)
	List(c *gin.Context, req *alertModels.SysAlertSinkTemplateReq) (*alertModels.SysAlertSinkTemplateListData, error)
	Remove(c *gin.Context, ids []string) error
}
