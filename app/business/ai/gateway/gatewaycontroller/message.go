package gatewaycontroller

import (
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// Message 智能体消息控制器。
type Message struct {
	service gatewayservice.IAIMessageService
}

// NewMessage 创建智能体消息控制器。
func NewMessage(service gatewayservice.IAIMessageService) *Message {
	return &Message{service: service}
}

// PrivateRoutes 注册智能体消息私有路由。
func (message *Message) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/ai/agent/messages")
	group.GET("/list", middlewares.HasPermission("ai:agent:messages:list"), message.List)
}

// PublicRoutes 注册智能体消息公开路由。
func (*Message) PublicRoutes(router *gin.RouterGroup) {

}

// List 根据会话ID读取消息列表
// @Summary 根据会话ID读取消息列表
// @Description 根据会话ID查询该会话下的全部消息记录，按创建时间正序返回
// @Tags 工业智能体/消息管理
// @Param object query gatewaymodels.MessageListReq true "消息列表查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} gatewaymodels.AIAgentMessageListData "获取成功"
// @Router /ai/agent/messages/list [get]
func (message *Message) List(c *gin.Context) {
	request := new(gatewaymodels.MessageListReq)
	err := c.BindQuery(request)
	if err != nil {
		baizeContext.ParameterError(c)
		zap.L().Error("query param error", zap.Error(err))
		return
	}
	if request.ConversationId == 0 {
		baizeContext.Waring(c, "消息id不能为空")
		return
	}

	data, err := message.service.GetMessage(c, request.ConversationId)
	if err != nil {
		zap.L().Error("get chat message error", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
