package aidatasetserviceimpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
)

type IAiPredictionExceptionServiceImpl struct {
	dao aidatasetdao.IAiPredictionExceptionDao
}

func NewIAiPredictionExceptionServiceImpl(dao aidatasetdao.IAiPredictionExceptionDao) aidatasetservice.IAiPredictionExceptionService {
	return &IAiPredictionExceptionServiceImpl{
		dao: dao,
	}
}

func (i *IAiPredictionExceptionServiceImpl) Set(c *gin.Context, data *aidatasetmodels.SetSysAiPredictionException) (*aidatasetmodels.SysAiPredictionException, error) {
	return i.dao.Set(c, data)
}
func (i *IAiPredictionExceptionServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *IAiPredictionExceptionServiceImpl) List(c *gin.Context, req *aidatasetmodels.SysAiPredictionExceptionListReq) (*aidatasetmodels.SysAiPredictionExceptionList, error) {
	return i.dao.List(c, req)
}
