package gatewaydao

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

// IAIAgentDao 智能体数据访问接口
type IAIAgentDao interface {
	Create(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error)
	Update(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.AIAgent, error)
	List(c *gin.Context, req *gatewaymodels.AIAgentQuery) (*gatewaymodels.AIAgentListData, error)
}
