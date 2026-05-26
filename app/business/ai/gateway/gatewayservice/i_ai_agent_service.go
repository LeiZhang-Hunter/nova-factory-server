package gatewayservice

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

// IAIAgentService 智能体配置服务接口
type IAIAgentService interface {
	Create(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error)
	Update(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.AIAgent, error)
	GetEnabledByType(c *gin.Context, agentType string) (*gatewaymodels.AIAgent, error)
	List(c *gin.Context, req *gatewaymodels.AIAgentQuery) (*gatewaymodels.AIAgentListData, error)
	// UpdateConfigVersion 更新版本
	UpdateConfigVersion(c *gin.Context, id int64, version string) error
}
