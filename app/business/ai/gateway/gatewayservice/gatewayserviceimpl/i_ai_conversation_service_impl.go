package gatewayserviceimpl

import (
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"

	"github.com/gin-gonic/gin"
)

type IAiConversationServiceImpl struct {
	dao        aidatasetdao.IAiConversationDao
	messageDao gatewaydao.IAIAgentMessageDao
}

func NewIAiConversationServiceImpl(dao aidatasetdao.IAiConversationDao, messageDao gatewaydao.IAIAgentMessageDao) gatewayservice.IAiConversationService {
	return &IAiConversationServiceImpl{
		dao:        dao,
		messageDao: messageDao,
	}
}

func (i *IAiConversationServiceImpl) Create(c *gin.Context, req *aidatasetmodels.SetAiConversation) (*aidatasetmodels.AiConversation, error) {
	return i.dao.Create(c, req)
}

func (i *IAiConversationServiceImpl) Update(c *gin.Context, req *aidatasetmodels.SetAiConversation) (*aidatasetmodels.AiConversation, error) {
	return i.dao.Update(c, req)
}

func (i *IAiConversationServiceImpl) List(c *gin.Context, req *aidatasetmodels.AiConversationQuery) (*aidatasetmodels.AiConversationListData, error) {
	conversationList, err := i.dao.List(c, req)
	if err != nil {
		zap.L().Error("conversation list failed", zap.Error(err))
		return conversationList, err
	}
	conversationIDs := make([]int64, 0, len(conversationList.Rows))
	conversationMap := make(map[int64]*aidatasetmodels.AiConversation, len(conversationList.Rows))
	for _, row := range conversationList.Rows {
		if row == nil || row.ID == 0 {
			continue
		}
		conversationIDs = append(conversationIDs, row.ID)
		conversationMap[row.ID] = row
	}
	if len(conversationIDs) == 0 {
		return conversationList, nil
	}
	messageList, messageErr := i.messageDao.List(c, &gatewaymodels.AIAgentMessageQuery{
		ConversationIDs: conversationIDs,
	})
	if messageErr != nil {
		zap.L().Error("conversation messages list failed", zap.Int("conversation_count", len(conversationIDs)), zap.Error(messageErr))
		return nil, messageErr
	}
	messageMap := make(map[int64][]*gatewaymodels.AIAgentMessage, len(conversationIDs))
	for _, message := range messageList.Rows {
		if message == nil || message.ConversationID == 0 {
			continue
		}
		rows := messageMap[message.ConversationID]
		if len(rows) >= 100 {
			continue
		}
		messageMap[message.ConversationID] = append(rows, message)
	}
	for conversationID, messages := range messageMap {
		row, exists := conversationMap[conversationID]
		if !exists {
			continue
		}
		row.Messages = reverseAIAgentMessages(messages)
		if row.Messages == nil {
			row.Messages = make([]*gatewaymodels.AIAgentMessage, 0)
		}
	}
	for k, _ := range conversationList.Rows {
		if conversationList.Rows[k].Messages == nil {
			conversationList.Rows[k].Messages = make([]*gatewaymodels.AIAgentMessage, 0)
		}
	}
	return conversationList, nil
}

func (i *IAiConversationServiceImpl) Remove(c *gin.Context, ids []int64) error {
	return i.dao.Remove(c, ids)
}

func reverseAIAgentMessages(rows []*gatewaymodels.AIAgentMessage) []*gatewaymodels.AIAgentMessage {
	if len(rows) <= 1 {
		return rows
	}
	result := make([]*gatewaymodels.AIAgentMessage, len(rows))
	for index := range rows {
		result[len(rows)-1-index] = rows[index]
	}
	return result
}
