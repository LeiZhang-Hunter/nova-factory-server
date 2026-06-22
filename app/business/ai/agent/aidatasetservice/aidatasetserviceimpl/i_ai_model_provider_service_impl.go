package aidatasetserviceimpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/utils/llm/embedding"
	"nova-factory-server/app/utils/store/embedder"
)

type IAiModelProviderServiceImpl struct {
	dao        aidatasetdao.IAiModelProviderDao
	settingDao aidatasetdao.IAiLLMSettingDao
	userLlmDao aidatasetdao.IAiUserLLMDao
}

func NewIAiModelProviderServiceImpl(dao aidatasetdao.IAiModelProviderDao,
	settingDao aidatasetdao.IAiLLMSettingDao,
	userLlmDao aidatasetdao.IAiUserLLMDao) aidatasetservice.IAiModelProviderService {
	i := &IAiModelProviderServiceImpl{
		dao:        dao,
		settingDao: settingDao,
		userLlmDao: userLlmDao,
	}
	embedder.RegisterStore(i)
	return i
}

func (i *IAiModelProviderServiceImpl) ListWithLLM(c *gin.Context, req *aidatasetmodels.SysAiModelProviderListReq) (*aidatasetmodels.SysAiModelProviderListData, error) {
	return i.dao.ListWithLLM(c, req)
}

func (i *IAiModelProviderServiceImpl) EmbeddingWithLLM(c *gin.Context) (embedder.EmbedderLlm, error) {
	info, err := i.settingDao.Get(c, &aidatasetmodels.GetSysAiLLMSettingReq{
		ID: 1,
	})
	if err != nil {
		return nil, err
	}

	if info == nil {
		return nil, nil
	}

	if info.EmbdID == "" {
		return nil, nil
	}

	providerID, modelId, err := embedding.ParseProviderModelID(info.EmbdID)
	if err != nil {
		return nil, err
	}

	llm, err := i.userLlmDao.GetByFidAndLlm(providerID, modelId)
	if err != nil {
		return nil, err
	}

	return llm, nil
}
