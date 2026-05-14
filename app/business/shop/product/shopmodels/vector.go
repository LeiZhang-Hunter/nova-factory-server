package shopmodels

type EmbeddingConfig struct {
	ProviderType string `json:"provider_type"`
	ProviderID   string `json:"provider_id"`
	APIEndpoint  string `json:"api_endpoint"`
	ModelID      string `json:"model_id"`
}

// GenVectorReq 生成向量req
type GenVectorReq struct {
	ID        int64            `json:"id"`
	Embedding *EmbeddingConfig `json:"embedding"`
}
