package gatewaycontroller

import (
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Conversations struct {
	service aidatasetservice.IAiConversationService
}

func NewConversations(service aidatasetservice.IAiConversationService) *Conversations {
	return &Conversations{service: service}
}

func (c *Conversations) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/gateway/conversations")
	group.GET("/list", middlewares.HasPermission("ai:gateway:conversations:list"), c.List)
	group.POST("/set", middlewares.HasPermission("ai:gateway:conversations:set"), c.Set)
	group.DELETE("/remove/:ids", middlewares.HasPermission("ai:gateway:conversations:remove"), c.Remove)
}

// List 查询会话列表
// @Summary 查询会话列表
// @Description 查询会话列表
// @Tags 网关/会话管理
// @Param object query aidatasetmodels.AiConversationQuery true "会话列表查询参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "查询成功"
// @Router /ai/gateway/conversations/list [get]
func (c *Conversations) List(ctx *gin.Context) {
	req := new(aidatasetmodels.AiConversationQuery)
	if err := ctx.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	data, err := c.service.List(ctx, req)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// Set 新增或修改会话
// @Summary 新增或修改会话
// @Description 传入id时修改，不传id时新增
// @Tags 网关/会话管理
// @Param object body aidatasetmodels.SetAiConversation true "会话设置参数"
// @Produce application/json
// @Success 200 {object} response.ResponseData "操作成功"
// @Router /ai/gateway/conversations/set [post]
func (c *Conversations) Set(ctx *gin.Context) {
	req := new(aidatasetmodels.SetAiConversation)
	if err := ctx.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}
	var (
		data *aidatasetmodels.AiConversation
		err  error
	)
	if req.ID > 0 {
		data, err = c.service.Update(ctx, req)
	} else {
		data, err = c.service.Create(ctx, req)
	}
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.SuccessData(ctx, data)
}

// Remove 删除会话
// @Summary 删除会话
// @Description 删除会话（软删除）
// @Tags 网关/会话管理
// @Param ids path string true "会话ID，多个用逗号分隔"
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /ai/gateway/conversations/remove/{ids} [delete]
func (c *Conversations) Remove(ctx *gin.Context) {
	ids := baizeContext.ParamInt64Array(ctx, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(ctx)
		return
	}
	if err := c.service.Remove(ctx, ids); err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	baizeContext.Success(ctx)
}
