package chatmodel

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// Store 提供查询 AI 模块模型连接信息的能力。
type Store interface {
	GetConnection(c *gin.Context, provider, model string) (LlmConnection, error)
}

// emptyStore 是未注册实现时的兜底，返回明确错误。
type emptyStore struct{}

// NewEmptyStore 创建空实现，便于测试与未注册场景。
func NewEmptyStore() Store {
	return &emptyStore{}
}

func (e *emptyStore) GetConnection(c *gin.Context, provider, model string) (LlmConnection, error) {
	return nil, errors.New("模型连接 store 未注册")
}
