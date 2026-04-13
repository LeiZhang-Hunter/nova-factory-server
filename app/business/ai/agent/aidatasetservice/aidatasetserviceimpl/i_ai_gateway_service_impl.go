package aidatasetserviceimpl

import (
	"errors"
	"nova-factory-server/app/utils/baizeContext"
	"strings"

	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	coreclient "nova-factory-server/app/business/ai/core/client"
	conversationagent "nova-factory-server/app/business/ai/core/gateway/agent"
	gatewayapi "nova-factory-server/app/business/ai/core/gateway/api"
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"

	"github.com/gin-gonic/gin"
)

type IAIGatewayServiceImpl struct {
	gatewayService gatewayservice.IAIGatewayService
	config         *coreclient.Config
}

func NewAIGatewayService(gatewayService gatewayservice.IAIGatewayService) aidatasetservice.IAIGatewayService {
	return &IAIGatewayServiceImpl{
		gatewayService: gatewayService,
		config: &coreclient.Config{
			Algorithm:          coreclient.AlgorithmRoundRobin,
			APIKeyHeader:       "Authorization",
			APIKeyPrefix:       "Bearer ",
			AgentGatewayHeader: "X-Agent-Gateway",
		},
	}
}

func (i *IAIGatewayServiceImpl) Chat(c *gin.Context, req *aidatasetmodels.SendMessageInput) (*gatewayapi.ChatResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.Content == "" {
		return nil, errors.New("content不能为空")
	}
	if req.TabID == "" {
		return nil, errors.New("tab_id不能为空")
	}
	conversations, err := i.newConversationsClient(c)
	if err != nil {
		return nil, err
	}
	return conversations.Chat(c, &gatewayapi.SendMessageInput{
		ConversationID: req.ConversationID,
		Content:        req.Content,
		TabID:          req.TabID,
		UserID:         baizeContext.GetUserId(c),
	})
}

// StopGeneration 停止上游正在进行中的模型生成。
func (i *IAIGatewayServiceImpl) StopGeneration(c *gin.Context, req *aidatasetmodels.StopGenerationInput) (*gatewayapi.StopGenerationResponse, error) {
	if req == nil {
		return nil, errors.New("request is nil")
	}
	if req.ConversationID == 0 {
		return nil, errors.New("conversation_id不能为空")
	}
	if strings.TrimSpace(req.TabID) == "" {
		return nil, errors.New("tab_id不能为空")
	}
	conversations, err := i.newConversationsClient(c)
	if err != nil {
		return nil, err
	}
	return conversations.StopGeneration(c, &gatewayapi.StopGenerationInput{
		ConversationID: req.ConversationID,
		AgentGateway:   "",
		TabID:          req.TabID,
	})
}

// newConversationsClient 创建会话网关客户端并注入当前可用网关。
func (i *IAIGatewayServiceImpl) newConversationsClient(c *gin.Context) (gatewayapi.Conversations, error) {
	enabled := true
	listData, err := i.gatewayService.List(c, &gatewaymodels.AIGatewayQuery{
		Enabled: &enabled,
		Page:    1,
		Size:    10,
	})
	if err != nil {
		return nil, err
	}
	if listData == nil || len(listData.Rows) == 0 {
		return nil, errors.New("未配置可用网关")
	}
	endpoints := make([]coreclient.Endpoint, 0, len(listData.Rows))
	for _, row := range listData.Rows {
		if row == nil {
			continue
		}
		enabledValue := false
		if row.Enabled != nil {
			enabledValue = *row.Enabled
		}
		endpoints = append(endpoints, coreclient.Endpoint{
			Name:    row.Name,
			BaseURL: row.BaseURL,
			APIKey:  row.APIKey,
			Enabled: enabledValue,
		})
	}
	if len(endpoints) == 0 {
		return nil, errors.New("未配置可用网关")
	}
	cfg := *i.config
	cfg.Endpoints = endpoints
	cli, err := coreclient.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return conversationagent.NewConversations(cli), nil
}
