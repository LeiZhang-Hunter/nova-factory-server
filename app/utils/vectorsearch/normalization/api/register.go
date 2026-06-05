package api

import (
	"fmt"
	"strings"
	"sync"
)

type StepFactory func() Step

var (
	stepRegistryMu sync.RWMutex
	stepRegistry   = make(map[string]StepFactory)
)

// Register 注册一个 step 工厂。
func Register(factory StepFactory) {
	if factory == nil {
		panic("normalization step factory is nil")
	}
	step := factory()
	if step == nil {
		panic("normalization step is nil")
	}
	typeName := strings.TrimSpace(step.Type())
	if typeName == "" {
		panic("normalization step type is empty")
	}

	stepRegistryMu.Lock()
	defer stepRegistryMu.Unlock()
	if _, exists := stepRegistry[typeName]; exists {
		panic(fmt.Sprintf("normalization step type %q already registered", typeName))
	}
	stepRegistry[typeName] = factory
}

// NewStep 根据类型创建 step 实例。
func NewStep(typeName string) (Step, error) {
	typeName = strings.TrimSpace(typeName)
	if typeName == "" {
		return nil, fmt.Errorf("normalization step type is empty")
	}

	stepRegistryMu.RLock()
	factory, ok := stepRegistry[typeName]
	stepRegistryMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("normalization step type %q is not registered", typeName)
	}
	step := factory()
	if step == nil {
		return nil, fmt.Errorf("normalization step type %q created nil step", typeName)
	}
	return step, nil
}
