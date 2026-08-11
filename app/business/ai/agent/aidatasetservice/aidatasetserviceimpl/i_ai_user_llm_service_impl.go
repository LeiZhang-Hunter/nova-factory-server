package aidatasetserviceimpl

import (
	"errors"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	chatmodelstore "nova-factory-server/app/utils/store/chatmodel"
	"strings"

	"github.com/gin-gonic/gin"
)

type IAiUserLLMServiceImpl struct {
	dao aidatasetdao.IAiUserLLMDao
}

var _ chatmodelstore.Store = (*IAiUserLLMServiceImpl)(nil)

func NewIAiUserLLMServiceImpl(dao aidatasetdao.IAiUserLLMDao) aidatasetservice.IAiUserLLMService {
	i := &IAiUserLLMServiceImpl{
		dao: dao,
	}
	chatmodelstore.RegisterStore(i)
	return i
}

// GetConnection 查询供应商下指定模型的连接信息，供其他模块通过 store 层使用。
// 只读取 user_id=0 的全局配置；查无返回 (nil, nil)。
func (i *IAiUserLLMServiceImpl) GetConnection(c *gin.Context, provider, model string) (chatmodelstore.LlmConnection, error) {
	row, err := i.dao.GetByFidAndLlm(provider, model)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, nil
	}
	return &chatmodelstore.LlmConnectionData{
		APIType:   row.APIType,
		APIKey:    row.APIKey,
		APIBase:   row.APIBase,
		MaxTokens: row.MaxTokens,
	}, nil
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
