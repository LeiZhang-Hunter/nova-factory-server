package gatewaycontroller

import (
	"nova-factory-server/app/business/admin/system/systemdao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

// SubAgent 子智能体配置控制器。
type SubAgent struct {
	service     gatewayservice.IAISubAgentService
	dictDataDao systemdao.IDictDataDao
}

// NewSubAgent 创建子智能体配置控制器。
func NewSubAgent(service gatewayservice.IAISubAgentService, dictDataDao systemdao.IDictDataDao) *SubAgent {
	return &SubAgent{
		service:     service,
		dictDataDao: dictDataDao,
	}
}

// PrivateRoutes 注册子智能体配置路由。
func (s *SubAgent) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/sub/agent/config")
	group.GET("/list", middlewares.HasPermission("ai:sub:agent:config:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("ai:sub:agent:config:info"), s.Info)
	group.POST("/set", middlewares.HasPermission("ai:sub:agent:config:set"), s.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:sub:agent:config:remove"), s.Delete)
}

func (s *SubAgent) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	router.RegisterPermission("GET", "/ai/sub/agent/config/list", "ai:sub:agent:config:list")
	router.RegisterPermission("GET", "/ai/sub/agent/config/info/:id", "ai:sub:agent:config:info")
	router.RegisterPermission("POST", "/ai/sub/agent/config/set", "ai:sub:agent:config:set")
	router.RegisterPermission("DELETE", "/ai/sub/agent/config/remove/:ids", "ai:sub:agent:config:remove")
}

// List 获取子智能体配置列表
// @Summary 获取子智能体配置列表
// @Description 获取子智能体配置列表
// @Tags 工业智能体/子智能体配置
// @Param object query gatewaymodels.AISubAgentQuery true "子智能体配置查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/sub/agent/list [get]
func (s *SubAgent) List(c *gin.Context) {
	req := new(gatewaymodels.AISubAgentQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Info 获取子智能体配置详情
// @Summary 获取子智能体配置详情
// @Description 根据ID获取子智能体配置详情
// @Tags 工业智能体/子智能体配置
// @Param id path int true "子智能体ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/sub/agent/info/{id} [get]
func (s *SubAgent) Info(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 保存子智能体配置
// @Summary 保存子智能体配置
// @Description 保存子智能体配置，id为空时新增，不为空时修改
// @Tags 工业智能体/子智能体配置
// @Param object body gatewaymodels.AISubAgentUpsert true "子智能体配置保存参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /ai/sub/agent/set [post]
func (s *SubAgent) Set(c *gin.Context) {
	req := new(gatewaymodels.AISubAgentUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *gatewaymodels.AISubAgent
		err  error
	)
	if req.ID > 0 {
		data, err = s.service.Update(c, req)
	} else {
		data, err = s.service.Create(c, req)
	}
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除子智能体配置
// @Summary 删除子智能体配置
// @Description 根据ID删除子智能体配置
// @Tags 工业智能体/子智能体配置
// @Param ids path string true "子智能体ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/sub/agent/remove/{ids} [delete]
func (s *SubAgent) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
