package gatewaycontroller

import (
	"go.uber.org/zap"
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
	group.GET("/query/:id", middlewares.HasPermission("ai:gateway:query"), a.GetByID)
	group.POST("/set", middlewares.HasPermission("ai:gateway:set"), a.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:gateway:remove"), a.Delete)
}

// List 获取网关列表
// @Summary 获取网关列表
// @Description 获取AI网关配置列表
// @Tags 工业智能体/网关管理
// @Param object query gatewaymodels.AIGatewayQuery true "网关查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/gateway/list [get]
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

// GetByID 获取网关详情
// @Summary 获取网关详情
// @Description 根据ID获取AI网关配置详情
// @Tags 工业智能体/网关管理
// @Param id path int true "网关ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/gateway/query/{id} [get]
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

// Set 保存网关配置
// @Summary 保存网关配置
// @Description 保存AI网关配置，id为空时新增，不为空时修改
// @Tags 工业智能体/网关管理
// @Param object body gatewaymodels.AIGatewayUpsert true "网关保存参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "保存成功"
// @Router /ai/gateway/set [post]
func (a *AIGateway) Set(c *gin.Context) {
	req := new(gatewaymodels.AIGatewayUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		zap.L().Error("gateway set failed", zap.Error(err))
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

// Delete 删除网关配置
// @Summary 删除网关配置
// @Description 根据ID删除AI网关配置
// @Tags 工业智能体/网关管理
// @Param ids path string true "网关ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/gateway/remove/{ids} [delete]
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
