package integration

type Config interface {
	// GetStatus 返回集成功能的开关状态，nil 或 false 表示该集成处于关闭/未配置状态
	GetStatus() *bool

	// GetData 返回配置的原始 JSON 字符串，由各实现按需反序列化为自己的配置结构
	GetData() string

	// GetOverrideURL 返回可覆盖默认回调地址的自定义URL，为空则使用默认值
	GetOverrideURL() string

	// GetMetadata 返回扩展元数据键值对，用于传递配置之外的附加上下文信息
	GetMetadata() map[string]any
}
