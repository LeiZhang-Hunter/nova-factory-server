package aidatasetserviceimpl

import (
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"

	"github.com/gin-gonic/gin"
)

type IAiUserLLMServiceImpl struct {
	dao aidatasetdao.IAiUserLLMDao
}

func NewIAiUserLLMServiceImpl(dao aidatasetdao.IAiUserLLMDao) aidatasetservice.IAiUserLLMService {
	return &IAiUserLLMServiceImpl{
		dao: dao,
	}
}

func (i *IAiUserLLMServiceImpl) Set(c *gin.Context, req *aidatasetmodels.SetSysUserLLM) (*aidatasetmodels.SysUserLLM, error) {
	return i.dao.Set(c, req)
}

func (i *IAiUserLLMServiceImpl) Get(c *gin.Context, req *aidatasetmodels.GetSysUserLLMReq) ([]*aidatasetmodels.SysUserLLM, error) {
	return i.dao.Get(c, req)
}
