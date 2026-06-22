package embedder

type EmbedderLlm interface {
	GetUserID() int64
	GetLLMFactory() string
	GetModelType() string
	GetLLMName() string
	GetAPIType() string
	GetAPIKey() string
	GetAPIBase() string
	GetMaxTokens() int64
	GetUsedTokens() int64
	GetStatus() string
}

type EmbedderLlmData struct {
	UserID     int64  `gorm:"column:user_id;primaryKey" json:"user_id,string"`
	LLMFactory string `gorm:"column:llm_factory;primaryKey" json:"llm_factory"`
	ModelType  string `gorm:"column:model_type" json:"model_type"`
	LLMName    string `gorm:"column:llm_name;primaryKey" json:"llm_name"`
	APIType    string `gorm:"column:api_type" json:"api_type"`
	APIKey     string `gorm:"column:api_key" json:"api_key"`
	APIBase    string `gorm:"column:api_base" json:"api_base"`
	MaxTokens  int64  `gorm:"column:max_tokens" json:"max_tokens"`
	UsedTokens int64  `gorm:"column:used_tokens" json:"used_tokens"`
	Status     string `gorm:"column:status" json:"status"`
}

func (e *EmbedderLlmData) GetUserID() int64 {
	return e.UserID
}
func (e *EmbedderLlmData) GetLLMFactory() string {
	return e.LLMFactory
}
func (e *EmbedderLlmData) GetModelType() string {
	return e.ModelType
}
func (e *EmbedderLlmData) GetLLMName() string {
	return e.LLMName
}
func (e *EmbedderLlmData) GetAPIType() string {
	return e.APIType
}
func (e *EmbedderLlmData) GetAPIKey() string {
	return e.APIKey
}
func (e *EmbedderLlmData) GetAPIBase() string {
	return e.APIBase
}
func (e *EmbedderLlmData) GetMaxTokens() int64 {
	return e.MaxTokens
}
func (e *EmbedderLlmData) GetUsedTokens() int64 {
	return e.UsedTokens
}
func (e *EmbedderLlmData) GetStatus() string {
	return e.Status
}
