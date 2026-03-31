package aidatasetserviceimpl

import (
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"

	"github.com/gin-gonic/gin"
)

type IAiPredictionServiceImpl struct {
	dao aidatasetdao.IAiPredictionListDao
}

func NewIAiPredictionServiceImpl(dao aidatasetdao.IAiPredictionListDao) aidatasetservice.IAiPredictionService {
	return &IAiPredictionServiceImpl{
		dao: dao,
	}
}

func (i *IAiPredictionServiceImpl) Set(c *gin.Context, data *aidatasetmodels.SetSysAiPrediction) (*aidatasetmodels.SysAiPrediction, error) {
	return i.dao.Set(c, data)
}
func (i *IAiPredictionServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *IAiPredictionServiceImpl) List(c *gin.Context, req *aidatasetmodels.SysAiPredictionListReq) (*aidatasetmodels.SysAiPredictionList, error) {
	return i.dao.List(c, req)
}
