package aidatasetserviceimpl

import (
	"context"
	"errors"
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
	if req.ConversationID == 0 {
		return nil, errors.New("conversation_id不能为空")
	}
	if req.Content == "" {
		return nil, errors.New("content不能为空")
	}
	if req.TabID == "" {
		return nil, errors.New("tab_id不能为空")
	}
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
	i.config.Endpoints = endpoints
	cli, err := coreclient.NewClient(*i.config)
	if err != nil {
		return nil, err
	}
	conversations := conversationagent.NewConversations(cli)
	return conversations.Chat(context.Background(), &gatewayapi.SendMessageInput{
		ConversationID: req.ConversationID,
		AgentGateway:   "",
		Content:        req.Content,
		TabID:          req.TabID,
	})
}
