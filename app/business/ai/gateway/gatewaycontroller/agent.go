package gatewaycontroller

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Agent 智能体控制器。
type Agent struct {
	service gatewayservice.IAIAgentService
}

// NewAgent 创建智能体控制器。
func NewAgent(service gatewayservice.IAIAgentService) *Agent {
	return &Agent{service: service}
}

// PrivateRoutes 注册智能体配置路由。
func (agent *Agent) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/agent/config")
	group.GET("/list", middlewares.HasPermission("ai:agent:config:list"), agent.List)
	group.GET("/query/:id", middlewares.HasPermission("ai:agent:config:query"), agent.GetByID)
	group.POST("/set", middlewares.HasPermission("ai:agent:config:set"), agent.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:agent:config:remove"), agent.Delete)
}

// List 获取智能体列表
// @Summary 获取智能体列表
// @Description 获取智能体列表
// @Tags 工业智能体/智能体管理
// @Param object query gatewaymodels.AIAgentQuery true "智能体查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/list [get]
func (agent *Agent) List(c *gin.Context) {
	req := new(gatewaymodels.AIAgentQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := agent.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 获取智能体详情
// @Summary 获取智能体详情
// @Description 根据ID获取智能体详情
// @Tags 工业智能体/智能体管理
// @Param id path int true "智能体ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/query/{id} [get]
func (agent *Agent) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := agent.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 保存智能体
// @Summary 保存智能体
// @Description 保存智能体，id为空时新增，不为空时修改
// @Tags 工业智能体/智能体管理
// @Param object body gatewaymodels.AIAgentUpsert true "智能体保存参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /ai/agent/set [post]
func (agent *Agent) Set(c *gin.Context) {
	req := new(gatewaymodels.AIAgentUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *gatewaymodels.AIAgent
		err  error
	)
	if req.ID != 0 {
		data, err = agent.service.Update(c, req)
	} else {
		data, err = agent.service.Create(c, req)
	}
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除智能体
// @Summary 删除智能体
// @Description 删除智能体
// @Tags 工业智能体/智能体管理
// @Param ids path string true "智能体ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/agent/remove/{ids} [delete]
func (agent *Agent) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := agent.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
