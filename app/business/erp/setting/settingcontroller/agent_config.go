package settingcontroller

import (
	"nova-factory-server/app/business/erp/setting/settingmodels"
	"nova-factory-server/app/business/erp/setting/settingservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AgentConfig struct {
	service settingservice.IAgentConfigService
}

func NewAgentConfig(service settingservice.IAgentConfigService) *AgentConfig {
	return &AgentConfig{service: service}
}

func (a *AgentConfig) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/setting/sale/agent-config")
	group.GET("/list", middlewares.HasPermission("erp:setting:sale:agentConfig:list"), a.List)
	group.GET("/:id", middlewares.HasPermission("erp:setting:sale:agentConfig:query"), a.GetByID)
	group.POST("/set", middlewares.HasPermission("erp:setting:sale:agentConfig:set"), a.Set)
	group.DELETE("/:ids", middlewares.HasPermission("erp:setting:sale:agentConfig:remove"), a.Delete)
}

// List 销售订单智能体配置列表
// @Summary 销售订单智能体配置列表
// @Description 按条件分页查询销售订单智能体配置
// @Tags ERP/系统配置
// @Security BearerAuth
// @Param object query settingmodels.AgentConfigQuery true "销售订单智能体配置查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/setting/sale/agent-config/list [get]
func (a *AgentConfig) List(c *gin.Context) {
	req := new(settingmodels.AgentConfigQuery)
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

// GetByID 销售订单智能体配置详情
// @Summary 销售订单智能体配置详情
// @Description 根据ID查询销售订单智能体配置详情
// @Tags ERP/系统配置
// @Security BearerAuth
// @Param id path int true "配置ID"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /erp/setting/sale/agent-config/{id} [get]
func (a *AgentConfig) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := a.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Set 销售订单智能体配置设置
// @Summary 销售订单智能体配置设置
// @Description 新增或修改销售订单智能体配置
// @Tags ERP/系统配置
// @Security BearerAuth
// @Accept application/json
// @Param body body settingmodels.AgentConfigUpsert true "销售订单智能体配置参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "设置成功"
// @Router /erp/setting/sale/agent-config/set [post]
func (a *AgentConfig) Set(c *gin.Context) {
	req := new(settingmodels.AgentConfigUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.L().Error("param error", zap.Error(err))
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *settingmodels.AgentConfig
		err  error
	)
	info, err := a.service.GetByID(c, 1)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if info != nil {
		req.ID = info.ID
	}
	if req.ID > 0 {
		data, err = a.service.Update(c, req)
	} else {
		data, err = a.service.Create(c, req)
	}
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除销售订单智能体配置
// @Summary 删除销售订单智能体配置
// @Description 根据ID删除销售订单智能体配置，多个ID用逗号分隔
// @Tags ERP/系统配置
// @Security BearerAuth
// @Param ids path string true "配置ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /erp/setting/sale/agent-config/{ids} [delete]
func (a *AgentConfig) Delete(c *gin.Context) {
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
