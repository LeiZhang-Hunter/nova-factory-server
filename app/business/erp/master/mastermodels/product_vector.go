package mastermodels

type ProductEmbeddingConfig struct {
	ProviderType string `json:"provider_type"`
	ProviderID   string `json:"provider_id"`
	APIEndpoint  string `json:"api_endpoint"`
	ModelID      string `json:"model_id"`
	ApiKey       string `json:"api_key"`
}

// ProductGenVectorReq 生成产品向量请求。
type ProductGenVectorReq struct {
	ID        int64                   `json:"id"`
	Embedding *ProductEmbeddingConfig `json:"embedding"`
}

// ProductGenAllVectorReq 全量生成产品向量请求。
type ProductGenAllVectorReq struct {
	BatchSize int                     `json:"batchSize"`
	Embedding *ProductEmbeddingConfig `json:"embedding"`
}

// ProductVectorResult 产品向量写入结果。
type ProductVectorResult struct {
	ProductID  int64  `json:"productId,string"`
	Collection string `json:"collection"`
	Dimension  int    `json:"dimension"`
}

// ProductVectorSearchReq 产品向量检索请求。
type ProductVectorSearchReq struct {
	Query     string                  `json:"query" binding:"required"`
	Limit     int                     `json:"limit"`
	Embedding *ProductEmbeddingConfig `json:"embedding"`
}

// ProductVectorBatchSearchReq 产品批量向量检索请求。
type ProductVectorBatchSearchReq struct {
	Queries   []string                `json:"queries" binding:"required"`
	Limit     int                     `json:"limit"`
	Embedding *ProductEmbeddingConfig `json:"embedding"`
}

// ProductVectorSearchItem 产品向量检索结果。
type ProductVectorSearchItem struct {
	ProductID     int64   `json:"productId,string"`
	Name          string  `json:"name"`
	BarCode       string  `json:"barCode"`
	CategoryId    int64   `json:"categoryId,string"`
	CategoryName  string  `json:"categoryName"`
	UnitId        int64   `json:"unitId,string"`
	UnitName      string  `json:"unitName"`
	Standard      string  `json:"standard"`
	Remark        string  `json:"remark"`
	ExpiryDay     int32   `json:"expiryDay"`
	Weight        float64 `json:"weight"`
	PurchasePrice float64 `json:"purchasePrice"`
	SalePrice     float64 `json:"salePrice"`
	MinPrice      float64 `json:"minPrice"`
	Content       string  `json:"content"`
	Score         float32 `json:"score"`
}

// ProductVectorSearchData 产品向量检索结果集。
type ProductVectorSearchData struct {
	Rows  []*ProductVectorSearchItem `json:"rows"`
	Total int64                      `json:"total"`
}

// ProductVectorBatchSearchItem 产品批量向量检索结果。
type ProductVectorBatchSearchItem struct {
	Query string                     `json:"query"`
	Rows  []*ProductVectorSearchItem `json:"rows"`
	Total int64                      `json:"total"`
}

// ProductVectorBatchSearchData 产品批量向量检索结果集。
type ProductVectorBatchSearchData struct {
	Rows  []*ProductVectorBatchSearchItem `json:"rows"`
	Total int64                           `json:"total"`
}

// ProductVectorTaskData 产品向量任务数据。
type ProductVectorTaskData struct {
	TaskID      string `json:"taskId"`
	ProgressKey string `json:"progressKey"`
	Status      string `json:"status"`
}

// ProductVectorTaskProgress 产品向量任务进度。
type ProductVectorTaskProgress struct {
	TaskID      string `json:"taskId"`
	Status      string `json:"status"`
	Total       int64  `json:"total"`
	Processed   int64  `json:"processed"`
	Success     int64  `json:"success"`
	Failed      int64  `json:"failed"`
	Progress    int    `json:"progress"`
	Message     string `json:"message"`
	CurrentID   int64  `json:"currentId,string"`
	CurrentName string `json:"currentName"`
	OperatorID  int64  `json:"operatorId,string"`
	Operator    string `json:"operator"`
	StartedAt   string `json:"startedAt"`
	UpdatedAt   string `json:"updatedAt"`
	FinishedAt  string `json:"finishedAt,omitempty"`
}

// ProductVectorUpsertItem 产品向量写入项。
type ProductVectorUpsertItem struct {
	Content string
	Vector  []float32
}
