package aidatasetserviceimpl

import (
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"

	"github.com/gin-gonic/gin"
)

type IAiConversationServiceImpl struct {
	dao aidatasetdao.IAiConversationDao
}

func NewIAiConversationServiceImpl(dao aidatasetdao.IAiConversationDao) aidatasetservice.IAiConversationService {
	return &IAiConversationServiceImpl{
		dao: dao,
	}
}

func (i *IAiConversationServiceImpl) Create(c *gin.Context, req *aidatasetmodels.SetAiConversation) (*aidatasetmodels.AiConversation, error) {
	return i.dao.Create(c, req)
}

func (i *IAiConversationServiceImpl) Update(c *gin.Context, req *aidatasetmodels.SetAiConversation) (*aidatasetmodels.AiConversation, error) {
	return i.dao.Update(c, req)
}

func (i *IAiConversationServiceImpl) List(c *gin.Context, req *aidatasetmodels.AiConversationQuery) (*aidatasetmodels.AiConversationListData, error) {
	return i.dao.List(c, req)
}

func (i *IAiConversationServiceImpl) Remove(c *gin.Context, ids []int64) error {
	return i.dao.Remove(c, ids)
}
