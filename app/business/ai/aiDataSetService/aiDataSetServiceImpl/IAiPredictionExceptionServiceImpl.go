package aiDataSetServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
)

type IAiPredictionExceptionServiceImpl struct {
	dao aiDataSetDao.IAiPredictionExceptionDao
}

func NewIAiPredictionExceptionServiceImpl(dao aiDataSetDao.IAiPredictionExceptionDao) aiDataSetService.IAiPredictionExceptionService {
	return &IAiPredictionExceptionServiceImpl{
		dao: dao,
	}
}

func (i *IAiPredictionExceptionServiceImpl) Set(c *gin.Context, data *aiDataSetModels.SetSysAiPredictionException) (*aiDataSetModels.SysAiPredictionException, error) {
	return i.dao.Set(c, data)
}
func (i *IAiPredictionExceptionServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *IAiPredictionExceptionServiceImpl) List(c *gin.Context, req *aiDataSetModels.SysAiPredictionExceptionListReq) (*aiDataSetModels.SysAiPredictionExceptionList, error) {
	return i.dao.List(c, req)
}
