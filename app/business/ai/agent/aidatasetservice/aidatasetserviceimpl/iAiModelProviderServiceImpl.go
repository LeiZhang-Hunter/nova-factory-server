package aidatasetserviceimpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
)

type IAiModelProviderServiceImpl struct {
	dao aidatasetdao.IAiModelProviderDao
}

func NewIAiModelProviderServiceImpl(dao aidatasetdao.IAiModelProviderDao) aidatasetservice.IAiModelProviderService {
	return &IAiModelProviderServiceImpl{
		dao: dao,
	}
}

func (i *IAiModelProviderServiceImpl) ListWithLLM(c *gin.Context, req *aidatasetmodels.SysAiModelProviderListReq) (*aidatasetmodels.SysAiModelProviderListData, error) {
	return i.dao.ListWithLLM(c, req)
}
