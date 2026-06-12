// 定义集成配置相关的接口与类型。
// Config 接口为各第三方系统的集成配置提供了统一的数据访问方式，
// 包括开关状态、原始配置数据、覆盖URL及扩展元数据。
package config

// Config 集成配置接口，定义第三方系统接入所需的通用配置访问方法。
// 每个集成系统（如管家婆、金蝶）的配置结构不同，但都需实现此接口，
// 以便上层通知器与同步逻辑能统一读取配置。
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
