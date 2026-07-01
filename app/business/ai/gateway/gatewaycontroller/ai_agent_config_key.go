package gatewaycontroller

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// AgentConfigKey API Key 控制器。
type AgentConfigKey struct {
	service gatewayservice.IAgentConfigKeyService
}

// NewAgentConfigKey 创建 API Key 控制器。
func NewAgentConfigKey(service gatewayservice.IAgentConfigKeyService) *AgentConfigKey {
	return &AgentConfigKey{service: service}
}

// PrivateRoutes 注册 API Key 路由。
func (a *AgentConfigKey) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/agent/config/key")
	group.GET("/list", middlewares.HasPermission("ai:agent:config:key:list"), a.List)
	group.POST("/set", middlewares.HasPermission("ai:agent:config:key:set"), a.Set)
	group.POST("/generate", middlewares.HasPermission("ai:agent:config:key:generate"), a.Generate)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:agent:config:key:remove"), a.Delete)
}

func (a *AgentConfigKey) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/ai/agent/config/key/list", "ai:agent:config:key:list")
	router.RegisterPermission("POST", "/ai/agent/config/key/set", "ai:agent:config:key:set")
	router.RegisterPermission("POST", "/ai/agent/config/key/generate", "ai:agent:config:key:generate")
	router.RegisterPermission("DELETE", "/ai/agent/config/key/remove/:ids", "ai:agent:config:key:remove")
}

// List 获取 API Key 列表
// @Summary 获取 API Key 列表
// @Description 获取 API Key 配置列表
// @Tags 工业智能体/ApiKey管理
// @Param object query gatewaymodels.AgentConfigKeyQuery true "API Key 查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/config/key/list [get]
func (a *AgentConfigKey) List(c *gin.Context) {
	req := new(gatewaymodels.AgentConfigKeyQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := a.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 新增 API Key
// @Summary 新增 API Key
// @Description 新增 API Key
// @Tags 工业智能体/ApiKey管理
// @Param object body gatewaymodels.AgentConfigKeyUpsert true "API Key 保存参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "新增成功"
// @Router /ai/agent/config/key/set [post]
func (a *AgentConfigKey) Set(c *gin.Context) {
	req := new(gatewaymodels.AgentConfigKeyUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := a.service.Create(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除 API Key
// @Summary 删除 API Key
// @Description 根据 ID 删除 API Key
// @Tags 工业智能体/ApiKey管理
// @Param ids path string true "API Key ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/agent/config/key/remove/{ids} [delete]
func (a *AgentConfigKey) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := a.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// Generate 生成 API Key
// @Summary 生成 API Key
// @Description 自动生成一个随机 API Key 并保存
// @Tags 工业智能体/ApiKey管理
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "生成成功"
// @Router /ai/agent/config/key/generate [post]
func (a *AgentConfigKey) Generate(c *gin.Context) {
	baizeContext.SuccessData(c, a.service.Generate())
}
