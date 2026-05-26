package gatewaycontroller

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// AgentOrchestration 智能体编排控制器。
type AgentOrchestration struct {
	service gatewayservice.IAIAgentOrchestrationService
}

// NewAgentOrchestration 创建智能体编排控制器。
func NewAgentOrchestration(service gatewayservice.IAIAgentOrchestrationService) *AgentOrchestration {
	return &AgentOrchestration{service: service}
}

// PrivateRoutes 注册智能体编排路由。
func (a *AgentOrchestration) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/agent/orchestration")
	group.GET("/info/:agentId", middlewares.HasPermission("ai:agent:orchestration:info"), a.Info)
	group.POST("/set", middlewares.HasPermission("ai:agent:orchestration:set"), a.Set)
	group.DELETE("/remove/:agentIds", middlewares.HasPermission("ai:agent:orchestration:remove"), a.Remove)
}

// Info 获取智能体编排详情
// @Summary 获取智能体编排详情
// @Description 根据智能体ID获取智能体编排详情
// @Tags 工业智能体/智能体编排
// @Param agentId path int true "智能体ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/orchestration/info/{agentId} [get]
func (a *AgentOrchestration) Info(c *gin.Context) {
	agentID := baizeContext.ParamInt64(c, "agentId")
	if agentID == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := a.service.Info(c, agentID)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 保存智能体编排配置
// @Summary 保存智能体编排配置
// @Description 根据智能体ID保存编排配置，存在时更新，不存在时新增
// @Tags 工业智能体/智能体编排
// @Param object body gatewaymodels.AIAgentOrchestrationUpsert true "智能体编排保存参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /ai/agent/orchestration/set [post]
func (a *AgentOrchestration) Set(c *gin.Context) {
	req := new(gatewaymodels.AIAgentOrchestrationUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := a.service.Set(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Remove 删除智能体编排配置
// @Summary 删除智能体编排配置
// @Description 根据智能体ID删除智能体编排配置
// @Tags 工业智能体/智能体编排
// @Param agentIds path string true "智能体ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/agent/orchestration/remove/{agentIds} [delete]
func (a *AgentOrchestration) Remove(c *gin.Context) {
	agentIDs := baizeContext.ParamInt64Array(c, "agentIds")
	if len(agentIDs) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := a.service.Remove(c, agentIDs); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
