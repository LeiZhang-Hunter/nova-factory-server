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

// GenAllVectorReq 全量生成向量请求
type GenAllVectorReq struct {
	TaskID    string           `json:"taskId"`    // 可选，指定任务ID续跑
	BatchSize int              `json:"batchSize"` // 每批处理商品数
	Embedding *EmbeddingConfig `json:"embedding"`
}

// GoodsVectorSearchReq 商品向量检索请求
type GoodsVectorSearchReq struct {
	Query                   string           `json:"query" binding:"required"` // 检索文本
	Limit                   int              `json:"limit"`                    // 返回条数
	Embedding               *EmbeddingConfig `json:"embedding"`                // 向量模型配置
	FallbackWithoutMetadata bool             `json:"fallbackWithoutMetadata"`  // metadata 过滤无结果时是否回退为无过滤检索
	SearchText              string           `json:"-"`
	IsSale                  *bool            `json:"-"` // 可选，按在售状态过滤
}

// GoodsVectorBatchSearchReq 商品批量向量检索请求
type GoodsVectorBatchSearchReq struct {
	Queries     []string         `json:"queries" binding:"required"` // 检索文本列表
	Limit       int              `json:"limit"`                      // 每条返回条数
	Embedding   *EmbeddingConfig `json:"embedding"`                  // 向量模型配置
	SearchTexts []string         `json:"-"`
	IsSale      *bool            `json:"-"` // 可选，按在售状态过滤
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

// GoodsVectorTaskData 商品向量任务数据
type GoodsVectorTaskData struct {
	TaskID      string `json:"taskId"`
	ProgressKey string `json:"progressKey"`
	Status      string `json:"status"`
}

// GoodsVectorTaskListData 商品向量任务列表
type GoodsVectorTaskListData struct {
	Rows  []*GoodsVectorTaskProgress `json:"rows"`
	Total int64                      `json:"total"`
}

// GoodsVectorTaskProgress 商品向量任务进度
type GoodsVectorTaskProgress struct {
	TaskID      string `json:"taskId"`
	Status      string `json:"status"`
	BatchSize   int    `json:"batchSize"`
	Total       int64  `json:"total"`
	Processed   int64  `json:"processed"`
	Success     int64  `json:"success"`
	Failed      int64  `json:"failed"`
	Progress    int    `json:"progress"`
	Message     string `json:"message"`
	CurrentID   int64  `json:"currentId"`
	CurrentName string `json:"currentName"`
	OperatorID  int64  `json:"operatorId"`
	Operator    string `json:"operator"`
	StartedAt   string `json:"startedAt"`
	UpdatedAt   string `json:"updatedAt"`
	FinishedAt  string `json:"finishedAt,omitempty"`
}

// GoodsVectorUpsertItem 商品向量写入项
type GoodsVectorUpsertItem struct {
	SkuID          int64
	SkuName        string
	SkuDescription string
	RetailPrice    float64
	Weight         float64
	Quantity       int64
	Metadata       map[string]any
	Content        string
	Vector         []float32
	IsSale         bool
}
