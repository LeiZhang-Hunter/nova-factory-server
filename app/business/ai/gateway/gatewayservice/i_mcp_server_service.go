package gatewayservice

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

type IMCPServerService interface {
	Create(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error)
	Update(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error)
	DeleteByIDs(c *gin.Context, ids []string) error
	List(c *gin.Context, req *gatewaymodels.MCPServerQuery) (*gatewaymodels.MCPServerListData, error)
}
