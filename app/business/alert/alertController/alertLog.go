package alertController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/alert/alertService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type AlertLog struct {
	service alertService.AlertLogService
}

func NewAlertLog(service alertService.AlertLogService) *AlertLog {
	return &AlertLog{
		service: service,
	}
}

func (log *AlertLog) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/alert/log")
	group.GET("/list", middlewares.HasPermission("alert:log:list"), log.List)
}

func (log *AlertLog) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/alert/log/v1")
	group.POST("/export", log.Export)
}

// Export 设置告警规则
// @Summary 设置告警规则
// @Description 设置告警规则
// @Tags 告警管理/告警规则管理
// @Param object body alertModels.AlertLogData true "助理列表参数"
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /api/alert/log/v1/export [post]
func (log *AlertLog) Export(c *gin.Context) {
	data := new(alertModels.AlertLogData)
	err := c.ShouldBindJSON(data)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	err = log.service.Export(c, *data)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

func (log *AlertLog) List(c *gin.Context) {

}
