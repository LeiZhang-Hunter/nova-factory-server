package gatewaydao

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

type IMCPServerDao interface {
	Create(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error)
	Update(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.MCPServer, error)
	List(c *gin.Context, req *gatewaymodels.MCPServerQuery) (*gatewaymodels.MCPServerListData, error)
}
