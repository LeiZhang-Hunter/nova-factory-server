package settingserviceimpl

import (
	"nova-factory-server/app/business/erp/setting/settingdao"
	"nova-factory-server/app/business/erp/setting/settingmodels"
	"nova-factory-server/app/business/erp/setting/settingservice"

	"github.com/gin-gonic/gin"
)

type AgentConfigServiceImpl struct {
	dao settingdao.IAgentConfigDao
}

func NewAgentConfigService(dao settingdao.IAgentConfigDao) settingservice.IAgentConfigService {
	return &AgentConfigServiceImpl{dao: dao}
}

func (a *AgentConfigServiceImpl) Create(c *gin.Context, req *settingmodels.AgentConfigUpsert) (*settingmodels.AgentConfig, error) {
	return a.dao.Create(c, req)
}

func (a *AgentConfigServiceImpl) Update(c *gin.Context, req *settingmodels.AgentConfigUpsert) (*settingmodels.AgentConfig, error) {
	return a.dao.Update(c, req)
}

func (a *AgentConfigServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return a.dao.DeleteByIDs(c, ids)
}

func (a *AgentConfigServiceImpl) GetByID(c *gin.Context, id int64) (*settingmodels.AgentConfig, error) {
	return a.dao.GetByID(c, id)
}

func (a *AgentConfigServiceImpl) List(c *gin.Context, req *settingmodels.AgentConfigQuery) (*settingmodels.AgentConfigListData, error) {
	return a.dao.List(c, req)
}
