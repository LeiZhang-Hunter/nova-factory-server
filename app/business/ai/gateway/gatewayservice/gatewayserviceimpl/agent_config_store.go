package gatewayserviceimpl

import (
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	agentconfigstore "nova-factory-server/app/utils/store/agentconfig"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// agentConfigStore 将 IAIAgentService 适配为共享 agentconfig.Store，
// 供 data 等其他模块在无业务包依赖的情况下查询已启用智能体配置。
type agentConfigStore struct {
	service gatewayservice.IAIAgentService
}

var _ agentconfigstore.Store = (*agentConfigStore)(nil)

func (s *agentConfigStore) GetEnabledByType(c *gin.Context, agentType string) (*agentconfigstore.AgentConfig, error) {
	agent, err := s.service.GetEnabledByType(c, agentType)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, nil
	}
	return &agentconfigstore.AgentConfig{
		Name:              agent.Name,
		Provider:          agent.DefaultLLMProviderID,
		Model:             agent.DefaultLLMModelID,
		Temperature:       agent.LLMTemperature,
		EnableTemperature: agentConfigBoolValue(agent.EnableLLMTemperature),
		TopP:              agent.LLMTopP,
		EnableTopP:        agentConfigBoolValue(agent.EnableLLMTopP),
		MaxTokens:         cast.ToInt(agent.LLMMaxTokens),
		EnableMaxTokens:   agentConfigBoolValue(agent.EnableLLMMaxTokens),
	}, nil
}

func agentConfigBoolValue(v *bool) bool {
	return v != nil && *v
}
