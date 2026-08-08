package dto

import "time"

// UpdateModelConfigRequest 更新数据平台模型配置请求。
type UpdateModelConfigRequest struct {
	Provider                string   `json:"provider" binding:"required,max=100"`
	Model                   string   `json:"model" binding:"required,max=200"`
	Temperature             *float64 `json:"temperature"`
	EnableTemperature       bool     `json:"enable_temperature"`
	TopP                    *float64 `json:"top_p"`
	EnableTopP              bool     `json:"enable_top_p"`
	MaxTokens               *int     `json:"max_tokens"`
	EnableMaxTokens         bool     `json:"enable_max_tokens"`
	MaxContextCount         int      `json:"max_context_count"`
	RetrievalTopK           int      `json:"retrieval_top_k"`
	RetrievalMatchThreshold *float64 `json:"retrieval_match_threshold"`
}

// ModelConfigResponse 模型配置响应。
type ModelConfigResponse struct {
	ID                      string    `json:"id"`
	Provider                string    `json:"provider"`
	Model                   string    `json:"model"`
	Temperature             float64   `json:"temperature"`
	EnableTemperature       bool      `json:"enable_temperature"`
	TopP                    float64   `json:"top_p"`
	EnableTopP              bool      `json:"enable_top_p"`
	MaxTokens               int       `json:"max_tokens"`
	EnableMaxTokens         bool      `json:"enable_max_tokens"`
	MaxContextCount         int       `json:"max_context_count"`
	RetrievalTopK           int       `json:"retrieval_top_k"`
	RetrievalMatchThreshold float64   `json:"retrieval_match_threshold"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
