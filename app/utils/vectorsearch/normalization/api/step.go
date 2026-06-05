package api

import "strings"

// Step 表示 pipeline 中的单个归一化步骤。
type Step interface {
	Name() string
	Type() string
	Apply(*Context) error
	StepInitializer
}

// StepInitializer 表示可在 pipeline 构建阶段按配置完成初始化的 step。
type StepInitializer interface {
	Init(config InterceptorConfig) error
}

// StepFunc 允许直接用函数构建 Step。
type StepFunc struct {
	name     string
	typeName string
	fn       func(*Context) error
}

// NewStepFunc 创建函数式 Step。
func NewStepFunc(name string, fn func(*Context) error) Step {
	return StepFunc{
		name:     strings.TrimSpace(name),
		typeName: "func",
		fn:       fn,
	}
}

// Name 返回步骤名。
func (s StepFunc) Name() string {
	if s.name == "" {
		return "anonymous"
	}
	return s.name
}

// Type 返回步骤类型。
func (s StepFunc) Type() string {
	if s.typeName == "" {
		return "func"
	}
	return s.typeName
}

// Apply 执行归一化逻辑。
func (s StepFunc) Apply(ctx *Context) error {
	if s.fn == nil {
		return nil
	}
	return s.fn(ctx)
}

// Init 对函数式 Step 不做额外初始化。
func (s StepFunc) Init(_ InterceptorConfig) error {
	return nil
}
