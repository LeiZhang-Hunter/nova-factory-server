package gatewaycontroller

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type MCPServer struct {
	service gatewayservice.IMCPServerService
}

func NewMCPServer(service gatewayservice.IMCPServerService) *MCPServer {
	return &MCPServer{service: service}
}

// PrivateRoutes 注册MCP服务配置路由
func (m *MCPServer) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/mcp/server")
	group.GET("/list", middlewares.HasPermission("ai:mcp:server:list"), m.List)
	group.POST("/set", middlewares.HasPermission("ai:mcp:server:set"), m.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:mcp:server:remove"), m.Delete)
}

// List 获取MCP服务配置列表
// @Summary 获取MCP服务配置列表
// @Description 获取MCP服务配置列表
// @Tags 工业智能体/MCP服务配置
// @Param object query gatewaymodels.MCPServerQuery true "MCP服务配置查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/mcp/server/list [get]
func (m *MCPServer) List(c *gin.Context) {
	req := new(gatewaymodels.MCPServerQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := m.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 保存MCP服务配置
// @Summary 保存MCP服务配置
// @Description 保存MCP服务配置，id为空时新增，不为空时修改
// @Tags 工业智能体/MCP服务配置
// @Param object body gatewaymodels.MCPServerUpsert true "MCP服务配置保存参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /ai/mcp/server/set [post]
func (m *MCPServer) Set(c *gin.Context) {
	req := new(gatewaymodels.MCPServerUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *gatewaymodels.MCPServer
		err  error
	)
	if req.ID != "" {
		data, err = m.service.Update(c, req)
	} else {
		data, err = m.service.Create(c, req)
	}
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除MCP服务配置
// @Summary 删除MCP服务配置
// @Description 删除MCP服务配置
// @Tags 工业智能体/MCP服务配置
// @Param ids path string true "MCP服务ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/mcp/server/remove/{ids} [delete]
func (m *MCPServer) Delete(c *gin.Context) {
	ids := baizeContext.ParamStringArray(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := m.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
