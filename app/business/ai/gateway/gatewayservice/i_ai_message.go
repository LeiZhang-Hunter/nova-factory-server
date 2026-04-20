package gatewayservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
)

type IAIMessageService interface {
	GetMessage(c *gin.Context, conversationId int64) (*gatewaymodels.AIAgentMessageListData, error)
}
