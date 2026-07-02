package gatewaycontroller

import (
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// AgentConfigPublishHistory 智能体配置发布历史控制器。
type AgentConfigPublishHistory struct {
	service gatewayservice.IAIAgentConfigPublishHistoryService
}

// NewAgentConfigPublishHistory 创建智能体配置发布历史控制器。
func NewAgentConfigPublishHistory(service gatewayservice.IAIAgentConfigPublishHistoryService) *AgentConfigPublishHistory {
	return &AgentConfigPublishHistory{service: service}
}

// PrivateRoutes 注册智能体配置发布历史路由。
func (a *AgentConfigPublishHistory) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/agent/config/publish/history")
	group.GET("/list", middlewares.HasPermission("ai:agent:config:publish:history:list"), a.List)
	group.GET("/info/:id", middlewares.HasPermission("ai:agent:config:publish:history:info"), a.Info)
	group.POST("/set", middlewares.HasPermission("ai:agent:config:publish:history:set"), a.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:agent:config:publish:history:remove"), a.Remove)
}

func (a *AgentConfigPublishHistory) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/ai/agent/config/publish/history/list", "ai:agent:config:publish:history:list")
	router.RegisterPermission("GET", "/ai/agent/config/publish/history/info/:id", "ai:agent:config:publish:history:info")
	router.RegisterPermission("POST", "/ai/agent/config/publish/history/set", "ai:agent:config:publish:history:set")
	router.RegisterPermission("DELETE", "/ai/agent/config/publish/history/remove/:ids", "ai:agent:config:publish:history:remove")
}

// List 获取智能体配置发布历史列表
// @Summary 获取智能体配置发布历史列表
// @Description 获取智能体配置发布历史列表
// @Tags 工业智能体/智能体配置发布历史
// @Param object query gatewaymodels.AIAgentConfigPublishHistoryQuery true "智能体配置发布历史查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/config/publish/history/list [get]
func (a *AgentConfigPublishHistory) List(c *gin.Context) {
	req := new(gatewaymodels.AIAgentConfigPublishHistoryQuery)
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

// Info 获取智能体配置发布历史详情
// @Summary 获取智能体配置发布历史详情
// @Description 根据ID获取智能体配置发布历史详情
// @Tags 工业智能体/智能体配置发布历史
// @Param id path int true "发布历史ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/config/publish/history/info/{id} [get]
func (a *AgentConfigPublishHistory) Info(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := a.service.Info(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 保存智能体配置发布历史
// @Summary 保存智能体配置发布历史
// @Description 保存智能体配置发布历史，id为空时新增，不为空时修改
// @Tags 工业智能体/智能体配置发布历史
// @Param object body gatewaymodels.AIAgentConfigPublishHistoryUpsert true "智能体配置发布历史保存参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /ai/agent/config/publish/history/set [post]
func (a *AgentConfigPublishHistory) Set(c *gin.Context) {
	req := new(gatewaymodels.AIAgentConfigPublishHistoryUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := a.service.Set(c, req)
	if err != nil {
		zap.L().Error("AgentConfigPublishHistory.Set failed", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Remove 删除智能体配置发布历史
// @Summary 删除智能体配置发布历史
// @Description 根据ID删除智能体配置发布历史
// @Tags 工业智能体/智能体配置发布历史
// @Param ids path string true "发布历史ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/agent/config/publish/history/remove/{ids} [delete]
func (a *AgentConfigPublishHistory) Remove(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := a.service.Remove(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
