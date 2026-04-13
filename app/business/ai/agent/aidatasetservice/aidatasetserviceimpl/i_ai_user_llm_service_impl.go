package aidatasetserviceimpl

import (
	"errors"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"strings"

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

func (i *IAiUserLLMServiceImpl) Remove(c *gin.Context, req *aidatasetmodels.GetSysUserLLMReq) error {
	if req == nil {
		return errors.New("参数不能为空")
	}
	req.LLMFactory = strings.TrimSpace(req.LLMFactory)
	if req.LLMFactory == "" {
		return errors.New("llm_factory不能为空")
	}
	return i.dao.Remove(c, req)
}
