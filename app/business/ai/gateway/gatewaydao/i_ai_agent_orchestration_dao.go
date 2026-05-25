package gatewaydao

import (
	"context"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

// IAIAgentOrchestrationDao 智能体编排数据访问接口。
type IAIAgentOrchestrationDao interface {
	Create(c *gin.Context, req *gatewaymodels.AIAgentOrchestrationUpsert) (*gatewaymodels.AIAgentOrchestration, error)
	UpdateByAgentID(c *gin.Context, req *gatewaymodels.AIAgentOrchestrationUpsert) (*gatewaymodels.AIAgentOrchestration, error)
	GetByAgentID(c *gin.Context, agentID int64) (*gatewaymodels.AIAgentOrchestration, error)
	DeleteByAgentIDs(c *gin.Context, agentIDs []int64) error
	GetConfigByAgentID(c context.Context, agentID int64) (*gatewaymodels.AIAgentOrchestration, error)
}
