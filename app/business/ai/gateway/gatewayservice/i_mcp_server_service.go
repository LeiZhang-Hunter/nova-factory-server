package gatewayservice

import (
	"context"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/utils/gin_mcp"

	"github.com/gin-gonic/gin"
)

type IMCPServerService interface {
	Create(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error)
	Update(c *gin.Context, req *gatewaymodels.MCPServerUpsert) (*gatewaymodels.MCPServer, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	List(c *gin.Context, req *gatewaymodels.MCPServerQuery) (*gatewaymodels.MCPServerListData, error)
	Probe(ctx context.Context, req *gatewaymodels.MCPServerProbeRequest) (*gatewaymodels.MCPServerProbeResult, error)
	ProbePerm(ctx *gin.Context, mcpServer *gin_mcp.GinMCP, req *gatewaymodels.MCPServerProbeRequest) (*gatewaymodels.MCPServerProbeResult, error)
}
