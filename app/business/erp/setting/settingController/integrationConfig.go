package settingController

import (
	"nova-factory-server/app/business/erp/setting/settingModels"
	"nova-factory-server/app/business/erp/setting/settingService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type IntegrationConfig struct {
	service settingService.IIntegrationConfigService
	host    string
}

func NewIntegrationConfig(service settingService.IIntegrationConfigService) *IntegrationConfig {
	host := viper.GetString("host")

	return &IntegrationConfig{
		service: service,
		host:    host,
	}
}

func (i *IntegrationConfig) PublicRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/setting/integration-config")
	group.GET("/oauth/callback", i.OAuthCallback)
	return
}

func (i *IntegrationConfig) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/setting/integration-config")
	group.GET("/list",
		middlewares.HasPermission("erp:setting:integrationConfig:list"), i.List)
	group.GET("/check-login-state",
		middlewares.HasPermission("erp:setting:integrationConfig:checkLoginState"), i.CheckLoginState)
	group.POST("/set",
		middlewares.HasPermission("erp:setting:integrationConfig:set"), i.Set)
}

// List 集成配置列表
// @Summary 集成配置列表
// @Description 按条件分页查询ERP集成配置
// @Tags ERP/系统配置
// @Security BearerAuth
// @Param type query string false "接入类型"
// @Param status query bool false "状态"
// @Param page query int false "页码"
// @Param size query int false "每页条数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/setting/integration-config/list [get]
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

// Set 集成配置设置
// @Summary 集成配置设置
// @Description 新增或修改ERP集成配置，仅允许存在一条启用配置
// @Tags ERP/系统配置
// @Security BearerAuth
// @Accept application/json
// @Param body body settingModels.IntegrationConfigSet true "集成配置参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/setting/integration-config/set [post]
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

// CheckLoginState 检查集成系统登录态
// @Summary 检查集成系统登录态
// @Description 按集成类型检查当前启用配置的登录态
// @Tags ERP/系统配置
// @Security BearerAuth
// @Param type query string true "接入类型"
// @Param checkUrl query string false "检查地址，优先级高于配置"
// @Param redirectUrl query string false "回调地址，优先级高于配置"
// @Produce application/json
// @Success 200 {object} response.ResponseData "检查成功"
// @Router /erp/setting/integration-config/check-login-state [get]
func (i *IntegrationConfig) CheckLoginState(c *gin.Context) {
	req := new(settingModels.IntegrationConfigCheckLoginReq)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.RedirectURL == "" {
		req.RedirectURL = i.host + "/erp/setting/integration-config/oauth/callback"
	}
	data, err := i.service.CheckLoginState(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// OAuthCallback 集成OAuth回调
// @Summary 集成OAuth回调
// @Description 接收管家婆OAuth回调的code和state
// @Tags ERP/系统配置
// @Param code query string true "授权码"
// @Param state query string false "透传字段"
// @Produce application/json
// @Success 200 {object} response.ResponseData "接收成功"
// @Router /erp/setting/integration-config/oauth/callback [get]
func (i *IntegrationConfig) OAuthCallback(c *gin.Context) {
	req := new(settingModels.IntegrationOAuthCallbackReq)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	baizeContext.SuccessData(c, &settingModels.IntegrationOAuthCallbackData{
		Code:  req.Code,
		State: req.State,
	})
}
