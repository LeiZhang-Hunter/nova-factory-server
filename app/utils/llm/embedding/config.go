package embedding

import "time"

// ProviderConfig 创建 Embedder 所需的配置
type ProviderConfig struct {
	// ProviderID 供应商 ID（如 chatwiki/openai）
	ProviderID string
	// ProviderType 供应商类型（openai, azure, ollama, gemini, anthropic）
	ProviderType string
	// APIKey 供应商的 API 密钥
	APIKey string
	// APIEndpoint 供应商 API 的基础 URL
	APIEndpoint string
	// ModelID 嵌入模型的 ID
	ModelID string
	// Dimension 向量维度（可选，某些模型支持）
	Dimension int
	// ExtraConfig 供应商特定的配置（JSON 格式）
	ExtraConfig string
	// Timeout 请求超时时间
	Timeout time.Duration
}
