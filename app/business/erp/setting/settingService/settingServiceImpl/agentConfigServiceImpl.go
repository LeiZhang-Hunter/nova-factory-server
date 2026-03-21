package settingServiceImpl

import (
	"nova-factory-server/app/business/erp/setting/settingDao"
	"nova-factory-server/app/business/erp/setting/settingModels"
	"nova-factory-server/app/business/erp/setting/settingService"

	"github.com/gin-gonic/gin"
)

type AgentConfigServiceImpl struct {
	dao settingDao.IAgentConfigDao
}

func NewAgentConfigService(dao settingDao.IAgentConfigDao) settingService.IAgentConfigService {
	return &AgentConfigServiceImpl{dao: dao}
}

func (a *AgentConfigServiceImpl) Create(c *gin.Context, req *settingModels.AgentConfigUpsert) (*settingModels.AgentConfig, error) {
	return a.dao.Create(c, req)
}

func (a *AgentConfigServiceImpl) Update(c *gin.Context, req *settingModels.AgentConfigUpsert) (*settingModels.AgentConfig, error) {
	return a.dao.Update(c, req)
}

func (a *AgentConfigServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return a.dao.DeleteByIDs(c, ids)
}

func (a *AgentConfigServiceImpl) GetByID(c *gin.Context, id int64) (*settingModels.AgentConfig, error) {
	return a.dao.GetByID(c, id)
}

func (a *AgentConfigServiceImpl) List(c *gin.Context, req *settingModels.AgentConfigQuery) (*settingModels.AgentConfigListData, error) {
	return a.dao.List(c, req)
}
