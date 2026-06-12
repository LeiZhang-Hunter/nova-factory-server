// 定义集成客户端的注册与工厂模式管理。
// 通过 Registry 统一管理各第三方系统的工厂函数，
// 支持按 Kind 动态创建 Service 实例，各适配器在 init() 中完成自注册。
package api

import (
	"errors"
	"nova-factory-server/app/utils/observer/integration/kind"
	"sync"
)

// Factory 创建 Service 客户端的工厂函数类型。
// 无参数并返回一个已初始化好的 Service 实例，
// 各适配器（如管家婆、金蝶）在 init() 中将自己的工厂注册到默认 Registry。
type Factory func() Service

// Registry 集成客户端工厂注册表，维护 Kind 到 Factory 的映射。
// 支持并发安全的注册与创建操作。
type Registry struct {
	mu        sync.RWMutex
	factories map[kind.Kind]Factory
}

// NewRegistry 创建新的注册表实例，内部初始化 factories map。
func NewRegistry() *Registry {
	return &Registry{
		factories: map[kind.Kind]Factory{},
	}
}

// Register 注册一个集成客户端的工厂函数。
// kind 和 factory 均不可为空，同一 kind 重复注册会覆盖。
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

// Create 根据集成类型创建对应的 Service 客户端实例。
// 每次调用都会通过工厂函数创建新实例。
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

// GetRegistry 获取全局默认注册表单例，线程安全。
func GetRegistry() *Registry {
	defaultRegistryOnce.Do(func() {
		defaultRegistry = NewRegistry()
	})
	return defaultRegistry
}

// Register 向全局默认注册表注册客户端工厂，各适配器在 init() 中调用。
func Register(kind kind.Kind, factory Factory) error {
	return GetRegistry().Register(kind, factory)
}
