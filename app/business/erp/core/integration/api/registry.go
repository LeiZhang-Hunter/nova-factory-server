package api

import (
	"errors"
	"sync"
)

type Factory func() Client

type Registry struct {
	mu        sync.RWMutex
	factories map[Kind]Factory
}

func NewRegistry() *Registry {
	return &Registry{
		factories: map[Kind]Factory{},
	}
}

func (r *Registry) Register(kind Kind, factory Factory) error {
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

func (r *Registry) Create(kind Kind) (Client, error) {
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

func GetRegistry() *Registry {
	defaultRegistryOnce.Do(func() {
		defaultRegistry = NewRegistry()
	})
	return defaultRegistry
}

func Register(kind Kind, factory Factory) error {
	return GetRegistry().Register(kind, factory)
}
