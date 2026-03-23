package aiDataSetController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
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
