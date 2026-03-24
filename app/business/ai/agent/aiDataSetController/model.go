package aiDataSetController

import (
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/business/ai/agent/aiDataSetService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Model struct {
	service aiDataSetService.IAiModelProviderService
}

func NewModel(service aiDataSetService.IAiModelProviderService) *Model {
	return &Model{
		service: service,
	}
}

func (m *Model) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/model")
	group.GET("/provider/list", middlewares.HasPermission("ai:model:provider:list"), m.ProviderList)
}

// ProviderList 模型供应商列表
// @Summary 模型供应商列表
// @Description 读取模型供应商及其下级LLM列表
// @Tags 工业智能体/模型配置
// @Param  object query aiDataSetModels.SysAiModelProviderListReq true "模型供应商列表参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/model/provider/list [get]
func (m *Model) ProviderList(c *gin.Context) {
	req := new(aiDataSetModels.SysAiModelProviderListReq)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	list, err := m.service.ListWithLLM(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list)
}
