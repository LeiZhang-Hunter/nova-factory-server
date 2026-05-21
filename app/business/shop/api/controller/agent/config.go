//go:build ai

package agent

import (
	"strings"

	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Config 商城智能体配置控制器。
type Config struct {
	service gatewayservice.IAIAgentService
}

// NewConfig 创建商城智能体配置控制器。
func NewConfig(service gatewayservice.IAIAgentService) *Config {
	return &Config{service: service}
}

// PrivateRoutes 注册商城智能体配置路由。
func (config *Config) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/agent/config")
	group.GET("/enabled", config.GetEnabledByType)
}

// GetEnabledByType 获取指定类型下已启用的智能体配置
// @Summary 获取指定类型下已启用的智能体配置
// @Description 小程序侧通过 agentType 查询已启用的智能体配置
// @Tags app接口/商城/App智能体配置
// @Param agentType query string true "智能体类型"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/agent/config/enabled [get]
func (config *Config) GetEnabledByType(c *gin.Context) {
	agentType := strings.TrimSpace(c.Query("agentType"))
	if agentType == "" {
		baizeContext.ParameterError(c)
		return
	}
	data, err := config.service.GetEnabledByType(c, agentType)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if data == nil {
		baizeContext.Waring(c, "智能体配置不存在")
		return
	}
	baizeContext.SuccessData(c, data)
}
