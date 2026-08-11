package chatmodel

import (
	"errors"
	"fmt"
	"strings"
	"time"

	agentconfigstore "nova-factory-server/app/utils/store/agentconfig"
	chatmodelstore "nova-factory-server/app/utils/store/chatmodel"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// defaultModelTimeout 模型请求超时，固定默认值。
const defaultModelTimeout = 120 * time.Second

// BuildAgentChatModelConfig 按智能体类型获取已启用智能体配置，
// 并合并模型连接信息，生成模型初始化配置，供各业务模块通用。
func BuildAgentChatModelConfig(c *gin.Context, agentType string) (*Config, error) {
	agentCfg, err := agentconfigstore.GetStore().GetEnabledByType(c, agentType)
	if err != nil {
		return nil, err
	}
	if agentCfg == nil || strings.TrimSpace(agentCfg.Provider) == "" || strings.TrimSpace(agentCfg.Model) == "" {
		return nil, fmt.Errorf("未配置已启用的智能体（type=%s），请先在「智能体配置」中创建并启用并设置默认模型", agentType)
	}
	conn, err := chatmodelstore.GetStore().GetConnection(c, agentCfg.Provider, agentCfg.Model)
	if err != nil {
		return nil, fmt.Errorf("读取模型连接信息失败: %w", err)
	}
	if conn == nil {
		return nil, fmt.Errorf("AI 模块未配置模型 %q 的连接信息（api_type/api_key）", agentCfg.Model)
	}
	return buildChatModelConfig(agentCfg, conn)
}

// buildChatModelConfig 合并智能体模型参数与 AI 连接信息，生成模型初始化配置。
func buildChatModelConfig(agentCfg *agentconfigstore.AgentConfig, conn chatmodelstore.LlmConnection) (*Config, error) {
	if agentCfg == nil || conn == nil {
		return nil, errors.New("智能体配置或连接信息不能为空")
	}
	protocol := strings.ToLower(strings.TrimSpace(conn.GetAPIType()))
	if protocol == "" {
		return nil, errors.New("模型连接缺少 api_type")
	}
	// max_tokens 仅在启用时下发；未启用则传 0（适配层会省略该参数，交由模型使用默认值）。
	// 注意：不能回退使用 ai_user_llm.max_tokens，该字段是模型上下文长度（可能高达百万），
	// 远超各家模型 max_tokens 输出上限，会触发 400。
	maxTokens := 0
	if agentCfg.EnableMaxTokens && agentCfg.MaxTokens > 0 {
		maxTokens = agentCfg.MaxTokens
	}
	cfg := &Config{
		Protocol:  protocol,
		APIKey:    conn.GetAPIKey(),
		BaseURL:   strings.TrimRight(strings.TrimSpace(conn.GetAPIBase()), "/"),
		Model:     strings.TrimSpace(agentCfg.Model),
		MaxTokens: maxTokens,
		Timeout:   defaultModelTimeout,
	}
	if agentCfg.EnableTemperature {
		temperature := cast.ToFloat32(agentCfg.Temperature)
		cfg.Temperature = &temperature
	}
	if agentCfg.EnableTopP {
		topP := cast.ToFloat32(agentCfg.TopP)
		cfg.TopP = &topP
	}
	return cfg, nil
}
