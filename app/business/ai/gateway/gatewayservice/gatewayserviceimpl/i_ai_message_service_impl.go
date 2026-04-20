package gatewayserviceimpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/gateway/gatewaydao"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
)

type IAiMessageServiceImpl struct {
	messageDao gatewaydao.IAIAgentMessageDao
}

func NewIAiMessageServiceImpl(messageDao gatewaydao.IAIAgentMessageDao) gatewayservice.IAIMessageService {
	return &IAiMessageServiceImpl{
		messageDao: messageDao,
	}
}

func (i *IAiMessageServiceImpl) GetMessage(c *gin.Context, conversationId int64) (*gatewaymodels.AIAgentMessageListData, error) {
	messageList, messageErr := i.messageDao.List(c, &gatewaymodels.AIAgentMessageQuery{
		ConversationIDs: []int64{
			conversationId,
		},
	})

	if messageErr != nil {
		return nil, messageErr
	}

	return messageList, nil
}
