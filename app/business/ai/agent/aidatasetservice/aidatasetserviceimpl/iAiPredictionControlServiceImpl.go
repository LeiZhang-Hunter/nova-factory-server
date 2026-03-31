package aidatasetserviceimpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
)

type IAiPredictionControlServiceImpl struct {
	dao aidatasetdao.IAiPredictionControlDao
}

func NewIAiPredictionControlServiceImpl(dao aidatasetdao.IAiPredictionControlDao) aidatasetservice.IAiPredictionControlService {
	return &IAiPredictionControlServiceImpl{
		dao: dao,
	}
}

func (i *IAiPredictionControlServiceImpl) Set(c *gin.Context, data *aidatasetmodels.SetSysAiPredictionControl) (*aidatasetmodels.SysAiPredictionControl, error) {
	return i.dao.Set(c, data)
}
func (i *IAiPredictionControlServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *IAiPredictionControlServiceImpl) List(c *gin.Context, req *aidatasetmodels.SysAiPredictionControlListReq) (*aidatasetmodels.SysAiPredictionControlList, error) {
	return i.dao.List(c, req)
}
