package settingController

import (
	"nova-factory-server/app/business/erp/setting/settingModels"
	"nova-factory-server/app/business/erp/setting/settingService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type IntegrationConfig struct {
	service settingService.IIntegrationConfigService
}

func NewIntegrationConfig(service settingService.IIntegrationConfigService) *IntegrationConfig {
	return &IntegrationConfig{service: service}
}

func (i *IntegrationConfig) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/setting/integration-config")
	group.GET("/list",
		middlewares.HasPermission("erp:setting:integrationConfig:list"), i.List)
	group.POST("/set",
		middlewares.HasPermission("erp:setting:integrationConfig:set"), i.Set)
}

func (i *IntegrationConfig) List(c *gin.Context) {
	req := new(settingModels.IntegrationConfigQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := i.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

func (i *IntegrationConfig) Set(c *gin.Context) {
	req := new(settingModels.IntegrationConfigSet)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	typeQuery := &settingModels.IntegrationConfigQuery{
		Type: req.Type,
		Page: 1,
		Size: 1,
	}
	typeData, err := i.service.List(c, typeQuery)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	willEnable := req.Status != nil && *req.Status
	if req.Status == nil && (typeData == nil || typeData.Total == 0) {
		willEnable = true
	}
	if willEnable {
		enabled := true
		enabledData, listErr := i.service.List(c, &settingModels.IntegrationConfigQuery{
			Status: &enabled,
			Page:   1,
			Size:   1000,
		})
		if listErr != nil {
			baizeContext.Waring(c, listErr.Error())
			return
		}
		for _, item := range enabledData.Rows {
			if item.Type != req.Type {
				baizeContext.Waring(c, "已经存在启用配置，请关闭掉启用配置")
				return
			}
		}
	}
	data, err := i.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
