package aiDataSetServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aiDataSetDao"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/business/ai/agent/aiDataSetService"
)

type IAiModelProviderServiceImpl struct {
	dao aiDataSetDao.IAiModelProviderDao
}

func NewIAiModelProviderServiceImpl(dao aiDataSetDao.IAiModelProviderDao) aiDataSetService.IAiModelProviderService {
	return &IAiModelProviderServiceImpl{
		dao: dao,
	}
}

func (i *IAiModelProviderServiceImpl) ListWithLLM(c *gin.Context, req *aiDataSetModels.SysAiModelProviderListReq) (*aiDataSetModels.SysAiModelProviderListData, error) {
	return i.dao.ListWithLLM(c, req)
}
