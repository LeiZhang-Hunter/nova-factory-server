package aidatasetdao

import (
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"

	"github.com/gin-gonic/gin"
)

type IAiConversationDao interface {
	Create(c *gin.Context, req *aidatasetmodels.SetAiConversation) (*aidatasetmodels.AiConversation, error)
	Update(c *gin.Context, req *aidatasetmodels.SetAiConversation) (*aidatasetmodels.AiConversation, error)
	List(c *gin.Context, req *aidatasetmodels.AiConversationQuery) (*aidatasetmodels.AiConversationListData, error)
	Remove(c *gin.Context, ids []int64) error
}
