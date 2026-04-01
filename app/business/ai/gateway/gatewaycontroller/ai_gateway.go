package gatewaycontroller

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type AIGateway struct {
	service gatewayservice.IAIGatewayService
}

func NewAIGateway(service gatewayservice.IAIGatewayService) *AIGateway {
	return &AIGateway{service: service}
}

func (a *AIGateway) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/gateway")
	group.GET("/list", middlewares.HasPermission("ai:gateway:list"), a.List)
	group.GET("/:id", middlewares.HasPermission("ai:gateway:query"), a.GetByID)
	group.POST("/set", middlewares.HasPermission("ai:gateway:set"), a.Set)
	group.DELETE("/:ids", middlewares.HasPermission("ai:gateway:remove"), a.Delete)
}

func (a *AIGateway) List(c *gin.Context) {
	req := new(gatewaymodels.AIGatewayQuery)
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

func (a *AIGateway) GetByID(c *gin.Context) {
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

func (a *AIGateway) Set(c *gin.Context) {
	req := new(gatewaymodels.AIGatewayUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	var (
		data *gatewaymodels.AIGateway
		err  error
	)
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

func (a *AIGateway) Delete(c *gin.Context) {
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
