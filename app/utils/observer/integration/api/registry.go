package api

import (
	"errors"
	"nova-factory-server/app/utils/observer/integration/kind"
	"sync"
)

// Factory 创建 Client 的工厂函数
type Factory func() Service

// Registry Client 工厂注册表
type Registry struct {
	mu        sync.RWMutex
	factories map[kind.Kind]Factory
}

// NewRegistry 创建注册表
func NewRegistry() *Registry {
	return &Registry{
		factories: map[kind.Kind]Factory{},
	}
}

// Register 注册客户端工厂
func (r *Registry) Register(kind kind.Kind, factory Factory) error {
	if kind == "" {
		return errors.New("kind不能为空")
	}
	if factory == nil {
		return errors.New("factory不能为空")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[kind] = factory
	return nil
}

// Create 按类型创建客户端
func (r *Registry) Create(kind kind.Kind) (Service, error) {
	if kind == "" {
		return nil, errors.New("kind不能为空")
	}
	r.mu.RLock()
	f, ok := r.factories[kind]
	r.mu.RUnlock()
	if !ok {
		return nil, errors.New("未支持的集成类型: " + string(kind))
	}
	return f(), nil
}

var (
	defaultRegistryOnce sync.Once
	defaultRegistry     *Registry
)

// GetRegistry 获取默认注册表单例
func GetRegistry() *Registry {
	defaultRegistryOnce.Do(func() {
		defaultRegistry = NewRegistry()
	})
	return defaultRegistry
}

// Register 向默认注册表注册客户端工厂
func Register(kind kind.Kind, factory Factory) error {
	return GetRegistry().Register(kind, factory)
}
