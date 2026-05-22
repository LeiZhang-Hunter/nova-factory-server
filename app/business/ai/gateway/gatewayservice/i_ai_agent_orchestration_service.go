package gatewayservice

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

// IAIAgentOrchestrationService 智能体编排服务接口。
type IAIAgentOrchestrationService interface {
	Set(c *gin.Context, req *gatewaymodels.AIAgentOrchestrationUpsert) (*gatewaymodels.AIAgentOrchestration, error)
	Info(c *gin.Context, agentID int64) (*gatewaymodels.AIAgentOrchestration, error)
	Remove(c *gin.Context, agentIDs []int64) error
}
