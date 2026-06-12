package settingcontroller

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"nova-factory-server/app/business/erp/setting/settingmodels"
	"nova-factory-server/app/business/erp/setting/settingservice"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type integrationLoginPayload struct {
	Type       string `json:"type"`
	State      string `json:"state"`
	Code       string `json:"code"`
	Token      string `json:"token"`
	ExpireDate string `json:"expireDate"`
	IssueDate  string `json:"issueDate"`
	AppKey     string `json:"appKey"`
	AppSecret  string `json:"appSecret"`
	ApiCode    int64  `json:"apiCode"`
	Message    string `json:"message"`
	UpdatedAt  string `json:"updatedAt"`
}

type IntegrationConfig struct {
	service settingservice.IIntegrationConfigService
	cache   cache.Cache
	host    string
}

func NewIntegrationConfig(service settingservice.IIntegrationConfigService, cache cache.Cache) *IntegrationConfig {
	host := viper.GetString("host")

	return &IntegrationConfig{
		service: service,
		host:    host,
		cache:   cache,
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
	req := new(settingmodels.IntegrationConfigQuery)
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
// @Param body body settingmodels.IntegrationConfigSet true "集成配置参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/setting/integration-config/set [post]
func (i *IntegrationConfig) Set(c *gin.Context) {
	req := new(settingmodels.IntegrationConfigSet)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	data, err := i.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// CheckLoginState 检查集成系统登录状态
// @Summary 检查集成系统登录状态
// @Description 按集成类型检查当前启用配置的登录状态
// @Tags ERP/系统配置
// @Security BearerAuth
// @Param type query string true "接入类型"
// @Param checkUrl query string false "检查地址，优先级高于配置"
// @Param redirectUrl query string false "回调地址，优先级高于配置"
// @Produce application/json
// @Success 200 {object} response.ResponseData "检查成功"
// @Router /erp/setting/integration-config/check-login-state [get]
func (i *IntegrationConfig) CheckLoginState(c *gin.Context) {
	req := new(settingmodels.IntegrationConfigCheckLoginReq)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.RedirectURL == "" {
		req.RedirectURL = i.host + "/erp/setting/integration-config/oauth/callback"
	}
	enableInfo, err := i.service.GetEnabled(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if enableInfo == nil {
		baizeContext.Success(c)
		return
	}

	service, err := enableInfo.Service()
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	if service == nil {
		baizeContext.Waring(c, "集成配置不存在")
		return
	}

	loginInfo, cacheErr := service.TokenGetter().GetTokenByCache(c, i.cache)
	if cacheErr == nil && loginInfo != nil && loginInfo.GetToken() != "" {
		baizeContext.SuccessData(c, loginInfo)
		return
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
	req := new(settingmodels.IntegrationOAuthCallbackReq)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}

	enableInfo, err := i.service.GetEnabled(c)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if enableInfo == nil {
		baizeContext.Success(c)
		return
	}

	service, err := enableInfo.Service()
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	if service == nil {
		baizeContext.Waring(c, "集成配置不存在")
		return
	}

	tokenData, err := service.TokenGetter().GetTokenByCode(c, enableInfo, req.Code)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}

	if tokenData.GetCode() != 0 {
		baizeContext.Waring(c, tokenData.GetMessage())
		return
	}

	if cacheErr := service.TokenGetter().SaveTokenToCache(c, i.cache, tokenData, 0); cacheErr != nil {
		baizeContext.Waring(c, cacheErr.Error())
		return
	}
	baizeContext.SuccessData(c, &settingmodels.IntegrationOAuthCallbackData{
		Code:       req.Code,
		State:      req.State,
		Token:      tokenData.GetToken(),
		ExpireDate: tokenData.GetExpireDate(),
		IssueDate:  tokenData.GetIssueDate(),
		AppKey:     tokenData.GetAppKey(),
		AppSecret:  tokenData.GetAppSecret(),
		Message:    tokenData.GetMessage(),
		ApiCode:    tokenData.GetCode(),
	})
}
