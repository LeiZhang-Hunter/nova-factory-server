package aidatasetservice

import (
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	gatewayapi "nova-factory-server/app/business/ai/core/gateway/api"

	"github.com/gin-gonic/gin"
)

type IAIGatewayService interface {
	Chat(c *gin.Context, req *aidatasetmodels.SendMessageInput) (*gatewayapi.ChatResponse, error)
}
