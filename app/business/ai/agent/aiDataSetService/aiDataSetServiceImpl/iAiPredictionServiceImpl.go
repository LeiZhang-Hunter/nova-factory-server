package aiDataSetServiceImpl

import (
	"nova-factory-server/app/business/ai/agent/aiDataSetDao"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/business/ai/agent/aiDataSetService"

	"github.com/gin-gonic/gin"
)

type IAiPredictionServiceImpl struct {
	dao aiDataSetDao.IAiPredictionListDao
}

func NewIAiPredictionServiceImpl(dao aiDataSetDao.IAiPredictionListDao) aiDataSetService.IAiPredictionService {
	return &IAiPredictionServiceImpl{
		dao: dao,
	}
}

func (i *IAiPredictionServiceImpl) Set(c *gin.Context, data *aiDataSetModels.SetSysAiPrediction) (*aiDataSetModels.SysAiPrediction, error) {
	return i.dao.Set(c, data)
}
func (i *IAiPredictionServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *IAiPredictionServiceImpl) List(c *gin.Context, req *aiDataSetModels.SysAiPredictionListReq) (*aiDataSetModels.SysAiPredictionList, error) {
	return i.dao.List(c, req)
}
