package normalization

import normalizationapi "nova-factory-server/app/utils/vectorsearch/normalization/api"

// Match 记录某个步骤触发的归一化命中信息。
type Match = normalizationapi.Match

// Result 表示 pipeline 执行后的归一化结果。
type Result = normalizationapi.Result

// Context 是 pipeline 在步骤间传递的上下文。
type Context = normalizationapi.Context

// NewContext 创建归一化上下文。
func NewContext(input string) *Context {
	return normalizationapi.NewContext(input)
}
