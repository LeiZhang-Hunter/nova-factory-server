package settingService

import (
	"nova-factory-server/app/business/erp/setting/settingModels"

	"github.com/gin-gonic/gin"
)

type IAgentConfigService interface {
	Create(c *gin.Context, req *settingModels.AgentConfigUpsert) (*settingModels.AgentConfig, error)
	Update(c *gin.Context, req *settingModels.AgentConfigUpsert) (*settingModels.AgentConfig, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*settingModels.AgentConfig, error)
	List(c *gin.Context, req *settingModels.AgentConfigQuery) (*settingModels.AgentConfigListData, error)
}
