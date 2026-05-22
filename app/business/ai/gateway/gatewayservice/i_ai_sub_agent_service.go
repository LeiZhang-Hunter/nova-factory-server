package gatewayservice

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

// IAISubAgentService 子智能体配置服务接口。
type IAISubAgentService interface {
	Create(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) (*gatewaymodels.AISubAgent, error)
	Update(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) (*gatewaymodels.AISubAgent, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.AISubAgent, error)
	List(c *gin.Context, req *gatewaymodels.AISubAgentQuery) (*gatewaymodels.AISubAgentListData, error)
	// ValidateType 校验子Agent
	ValidateType(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) error
}
