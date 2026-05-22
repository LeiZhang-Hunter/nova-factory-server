package gatewaydao

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

// IAISubAgentDao 子智能体配置数据访问接口。
type IAISubAgentDao interface {
	Create(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) (*gatewaymodels.AISubAgent, error)
	Update(c *gin.Context, req *gatewaymodels.AISubAgentUpsert) (*gatewaymodels.AISubAgent, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.AISubAgent, error)
	List(c *gin.Context, req *gatewaymodels.AISubAgentQuery) (*gatewaymodels.AISubAgentListData, error)
	GetByName(c *gin.Context, name string) (*gatewaymodels.AISubAgent, error)
}
