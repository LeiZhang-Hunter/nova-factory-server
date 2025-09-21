package aiDataSetServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
)

type IAiPredictionControlServiceImpl struct {
	dao aiDataSetDao.IAiPredictionControlDao
}

func NewIAiPredictionControlServiceImpl(dao aiDataSetDao.IAiPredictionControlDao) aiDataSetService.IAiPredictionControlService {
	return &IAiPredictionControlServiceImpl{
		dao: dao,
	}
}

func (i *IAiPredictionControlServiceImpl) Set(c *gin.Context, data *aiDataSetModels.SetSysAiPredictionControl) (*aiDataSetModels.SysAiPredictionControl, error) {
	return i.dao.Set(c, data)
}
func (i *IAiPredictionControlServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *IAiPredictionControlServiceImpl) List(c *gin.Context, req *aiDataSetModels.SysAiPredictionControlListReq) (*aiDataSetModels.SysAiPredictionControlList, error) {
	return i.dao.List(c, req)
}
