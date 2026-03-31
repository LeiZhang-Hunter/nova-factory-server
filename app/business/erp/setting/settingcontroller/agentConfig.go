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
