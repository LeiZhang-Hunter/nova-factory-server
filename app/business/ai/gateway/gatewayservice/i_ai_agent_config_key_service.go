package gatewayservice

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

// IAgentConfigKeyService API Key Service 接口。
type IAgentConfigKeyService interface {
	Create(c *gin.Context, req *gatewaymodels.AgentConfigKeyUpsert) (*gatewaymodels.AgentConfigKey, error)
	Update(c *gin.Context, req *gatewaymodels.AgentConfigKeyUpsert) (*gatewaymodels.AgentConfigKey, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.AgentConfigKey, error)
	List(c *gin.Context, req *gatewaymodels.AgentConfigKeyQuery) (*gatewaymodels.AgentConfigKeyListData, error)
	Generate() string
}
