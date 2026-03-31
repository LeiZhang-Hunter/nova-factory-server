package settingdao

import (
	"nova-factory-server/app/business/erp/setting/settingmodels"

	"github.com/gin-gonic/gin"
)

type IAgentConfigDao interface {
	Create(c *gin.Context, req *settingmodels.AgentConfigUpsert) (*settingmodels.AgentConfig, error)
	Update(c *gin.Context, req *settingmodels.AgentConfigUpsert) (*settingmodels.AgentConfig, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*settingmodels.AgentConfig, error)
	List(c *gin.Context, req *settingmodels.AgentConfigQuery) (*settingmodels.AgentConfigListData, error)
}
