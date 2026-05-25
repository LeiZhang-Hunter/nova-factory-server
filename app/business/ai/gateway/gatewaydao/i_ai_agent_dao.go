package gatewaydao

import (
	"context"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

// IAIAgentDao 智能体数据访问接口
type IAIAgentDao interface {
	Create(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error)
	Update(c *gin.Context, req *gatewaymodels.AIAgentUpsert) (*gatewaymodels.AIAgent, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.AIAgent, error)
	GetEnabledByType(c *gin.Context, agentType string) (*gatewaymodels.AIAgent, error)
	List(c *gin.Context, req *gatewaymodels.AIAgentQuery) (*gatewaymodels.AIAgentListData, error)
	UpdateConfigVersion(c *gin.Context, id int64, version string) error
	GetEnable(c context.Context) ([]*gatewaymodels.AIAgent, error)
}
