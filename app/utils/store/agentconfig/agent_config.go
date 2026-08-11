package agentconfig

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// AgentConfig 供其他模块使用的智能体模型参数快照。
// 由 AI 模块注册实现，避免业务模块之间直接依赖。
type AgentConfig struct {
	Name              string
	Provider          string
	Model             string
	Temperature       float64
	EnableTemperature bool
	TopP              float64
	EnableTopP        bool
	MaxTokens         int
	EnableMaxTokens   bool
}

// Store 提供按类型查询已启用智能体配置的能力。
type Store interface {
	GetEnabledByType(c *gin.Context, agentType string) (*AgentConfig, error)
}

// emptyStore 是未注册实现时的兜底，返回明确错误。
type emptyStore struct{}

// NewEmptyStore 创建空实现，便于测试与未注册场景。
func NewEmptyStore() Store {
	return &emptyStore{}
}

func (e *emptyStore) GetEnabledByType(c *gin.Context, agentType string) (*AgentConfig, error) {
	return nil, errors.New("智能体配置 store 未注册")
}
