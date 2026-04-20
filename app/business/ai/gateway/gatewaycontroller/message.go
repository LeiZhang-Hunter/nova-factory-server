package gatewaycontroller

import (
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
// @Description 根据会话ID读取消息列表
// @Tags 工业智能体/消息管理
// @Param conversationId query int true "会话ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /ai/agent/messages/list [get]
func (message *Message) List(c *gin.Context) {
	conversationID := baizeContext.QueryInt64(c, "conversationId")
	if conversationID == 0 {
		baizeContext.Waring(c, "消息id不能为空")
		return
	}

	data, err := message.service.GetMessage(c, conversationID)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
