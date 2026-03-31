package aiDataSetServiceImpl

import (
	"nova-factory-server/app/business/ai/agent/aiDataSetDao"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/business/ai/agent/aiDataSetService"

	"github.com/gin-gonic/gin"
)

type IAiUserLLMServiceImpl struct {
	dao aiDataSetDao.IAiUserLLMDao
}

func NewIAiUserLLMServiceImpl(dao aiDataSetDao.IAiUserLLMDao) aiDataSetService.IAiUserLLMService {
	return &IAiUserLLMServiceImpl{
		dao: dao,
	}
}

func (i *IAiUserLLMServiceImpl) Set(c *gin.Context, req *aiDataSetModels.SetSysUserLLM) (*aiDataSetModels.SysUserLLM, error) {
	return i.dao.Set(c, req)
}

func (i *IAiUserLLMServiceImpl) Get(c *gin.Context, req *aiDataSetModels.GetSysUserLLMReq) ([]*aiDataSetModels.SysUserLLM, error) {
	return i.dao.Get(c, req)
}
