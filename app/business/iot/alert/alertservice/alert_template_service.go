package alertservice

import (
	"nova-factory-server/app/business/iot/alert/alertmodels"

	"github.com/gin-gonic/gin"
)

type AlertTemplateService interface {
	Create(c *gin.Context, data *alertmodels.SetSysAlertSinkTemplate) (*alertmodels.SysAlertSinkTemplate, error)
	Update(c *gin.Context, data *alertmodels.SetSysAlertSinkTemplate) (*alertmodels.SysAlertSinkTemplate, error)
	List(c *gin.Context, req *alertmodels.SysAlertSinkTemplateReq) (*alertmodels.SysAlertSinkTemplateListData, error)
	Remove(c *gin.Context, ids []string) error
}
