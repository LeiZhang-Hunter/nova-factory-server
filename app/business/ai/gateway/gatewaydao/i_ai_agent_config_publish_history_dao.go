package gatewaydao

import (
	"context"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IAIAgentConfigPublishHistoryDao 智能体配置发布历史数据访问接口。
type IAIAgentConfigPublishHistoryDao interface {
	Create(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) (*gatewaymodels.AIAgentConfigPublishHistory, error)
	CreateWithTx(c *gin.Context, tx *gorm.DB, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) (*gatewaymodels.AIAgentConfigPublishHistory, error)
	Update(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) (*gatewaymodels.AIAgentConfigPublishHistory, error)
	UpdateWithTx(c *gin.Context, tx *gorm.DB, req *gatewaymodels.AIAgentConfigPublishHistoryUpsert) (*gatewaymodels.AIAgentConfigPublishHistory, error)
	GetByID(c *gin.Context, id int64) (*gatewaymodels.AIAgentConfigPublishHistory, error)
	GetByAgentIDAndVersion(c *gin.Context, agentID int64, version string) (*gatewaymodels.AIAgentConfigPublishHistory, error)
	GetByVersion(c context.Context, version string) (*gatewaymodels.AIAgentConfigPublishHistory, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	DeleteByIDsWithTx(c *gin.Context, tx *gorm.DB, ids []int64) error
	List(c *gin.Context, req *gatewaymodels.AIAgentConfigPublishHistoryQuery) (*gatewaymodels.AIAgentConfigPublishHistoryListData, error)
	// GetConfigsByAgentIdAndVersion 读取agent配置信息
	GetConfigsByAgentIdAndVersion(c context.Context, conditionMap map[int64]string) ([]*gatewaymodels.AIAgentConfigPublishHistory, error)
}
