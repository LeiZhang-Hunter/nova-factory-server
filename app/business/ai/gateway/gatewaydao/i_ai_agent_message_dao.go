package gatewaydao

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

type IAIAgentMessageDao interface {
	Create(c *gin.Context, req *gatewaymodels.AIAgentMessageUpsert) (*gatewaymodels.AIAgentMessage, error)
	Update(c *gin.Context, req *gatewaymodels.AIAgentMessageUpsert) (*gatewaymodels.AIAgentMessage, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.AIAgentMessage, error)
	List(c *gin.Context, req *gatewaymodels.AIAgentMessageQuery) (*gatewaymodels.AIAgentMessageListData, error)
}
