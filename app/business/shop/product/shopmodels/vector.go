package shopmodels

type EmbeddingConfig struct {
	ProviderType string `json:"provider_type"`
	ProviderID   string `json:"provider_id"`
	APIEndpoint  string `json:"api_endpoint"`
	ModelID      string `json:"model_id"`
	ApiKey       string `json:"api_key"`
}

// GenVectorReq 生成向量req
type GenVectorReq struct {
	ID        int64            `json:"id"`
	Embedding *EmbeddingConfig `json:"embedding"`
}

// GoodsVectorSearchReq 商品向量检索请求
type GoodsVectorSearchReq struct {
	Query     string           `json:"query" binding:"required"` // 检索文本
	Limit     int              `json:"limit"`                    // 返回条数
	Embedding *EmbeddingConfig `json:"embedding"`                // 向量模型配置
}

// GoodsVectorBatchSearchReq 商品批量向量检索请求
type GoodsVectorBatchSearchReq struct {
	Queries   []string         `json:"queries" binding:"required"` // 检索文本列表
	Limit     int              `json:"limit"`                      // 每条返回条数
	Embedding *EmbeddingConfig `json:"embedding"`                  // 向量模型配置
}

// GoodsVectorSearchItem 商品向量检索结果
type GoodsVectorSearchItem struct {
	GoodsDBID      int64   `json:"goodsDbId"`
	GoodsName      string  `json:"goodsName"`
	GoodsCode      string  `json:"goodsCode"`
	CategoryName   string  `json:"categoryName"`
	Description    string  `json:"description"`
	SkuID          int64   `json:"skuId"`
	SkuName        string  `json:"skuName"`
	SkuDescription string  `json:"skuDescription"`
	RetailPrice    float64 `json:"retailPrice"`
	Weight         float64 `json:"weight"`
	Quantity       int64   `json:"quantity"`
	Content        string  `json:"content"`
	Score          float32 `json:"score"`
}

// GoodsVectorSearchData 商品向量检索结果集
type GoodsVectorSearchData struct {
	Rows  []*GoodsVectorSearchItem `json:"rows"`
	Total int64                    `json:"total"`
}

// GoodsVectorBatchSearchItem 商品批量向量检索结果
type GoodsVectorBatchSearchItem struct {
	Query string                   `json:"query"`
	Rows  []*GoodsVectorSearchItem `json:"rows"`
	Total int64                    `json:"total"`
}

// GoodsVectorBatchSearchData 商品批量向量检索结果集
type GoodsVectorBatchSearchData struct {
	Rows  []*GoodsVectorBatchSearchItem `json:"rows"`
	Total int64                         `json:"total"`
}

// GoodsVectorUpsertItem 商品向量写入项
type GoodsVectorUpsertItem struct {
	SkuID          int64
	SkuName        string
	SkuDescription string
	RetailPrice    float64
	Weight         float64
	Quantity       int64
	Content        string
	Vector         []float32
}
