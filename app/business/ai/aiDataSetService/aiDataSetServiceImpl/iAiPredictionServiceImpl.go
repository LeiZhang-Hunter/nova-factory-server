package aiDataSetServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
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
