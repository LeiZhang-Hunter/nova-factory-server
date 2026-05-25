package gatewayservice

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	v1 "nova-factory-server/app/utils/grpc/confighotload/v1"

	"github.com/gin-gonic/gin"
)

// IAIAgentConfigPublisher 智能体配置发布器接口。
type IAIAgentConfigPublisher interface {
	BroadcastByGrpcClient(ctx *gin.Context, req *v1.AgentBroadcastRequest) (*v1.AgentBroadcastReply, error)
}

// IAIAgentConfigPublishHistoryService 智能体配置发布历史服务接口。
type IAIAgentConfigPublishHistoryService interface {
	Set(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) (*gatewaymodels.AIAgentConfigPublishHistory, error)
	Info(c *gin.Context, id int64) (*gatewaymodels.AIAgentConfigPublishHistory, error)
	List(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryQuery) (*gatewaymodels.AIAgentConfigPublishHistoryListData, error)
	Remove(c *gin.Context, ids []int64) error
}
