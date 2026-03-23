package aiDataSetModels

type ModelRequest struct {
	Model          string   `json:"model" binding:"required"`
	Prompt         string   `json:"prompt" binding:"required"`
	SystemPrompt   string   `json:"systemPrompt"`
	ProviderAPIKey string   `json:"providerApiKey"`
	ProviderURL    string   `json:"providerUrl"`
	MaxTokens      int      `json:"maxTokens"`
	Temperature    *float32 `json:"temperature"`
	TopP           *float32 `json:"topP"`
	TopK           *int32   `json:"topK"`
	TimeoutSeconds int      `json:"timeoutSeconds"`
}

type ModelResponse struct {
	Provider         string `json:"provider"`
	Model            string `json:"model"`
	Content          string `json:"content"`
	PromptTokens     int    `json:"promptTokens"`
	CompletionTokens int    `json:"completionTokens"`
	TotalTokens      int    `json:"totalTokens"`
}
