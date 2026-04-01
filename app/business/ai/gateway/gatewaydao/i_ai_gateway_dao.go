package gatewaydao

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

type IAIGatewayDao interface {
	Create(c *gin.Context, req *gatewaymodels.AIGatewayUpsert) (*gatewaymodels.AIGateway, error)
	Update(c *gin.Context, req *gatewaymodels.AIGatewayUpsert) (*gatewaymodels.AIGateway, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.AIGateway, error)
	List(c *gin.Context, req *gatewaymodels.AIGatewayQuery) (*gatewaymodels.AIGatewayListData, error)
}
