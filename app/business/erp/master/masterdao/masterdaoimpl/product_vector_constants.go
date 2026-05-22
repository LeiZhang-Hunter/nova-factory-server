package masterdaoimpl

import "nova-factory-server/app/business/erp/master/masterdao"

const (
	defaultProductVectorCollection = "erp_master_product_vectors"

	productTextBm25EmbFunctionName = "text_bm25_emb"

	productVectorPKField            = "product_id"
	productVectorNameField          = "name"
	productVectorBarCodeField       = "bar_code"
	productVectorCategoryIDField    = "category_id"
	productVectorCategoryNameField  = "category_name"
	productVectorUnitIDField        = "unit_id"
	productVectorUnitNameField      = "unit_name"
	productVectorStandardField      = "standard"
	productVectorRemarkField        = "remark"
	productVectorExpiryDayField     = "expiry_day"
	productVectorWeightField        = "weight"
	productVectorPurchasePriceField = "purchase_price"
	productVectorSalePriceField     = "sale_price"
	productVectorMinPriceField      = "min_price"
	productVectorContentField       = "content"
	productVectorEmbeddingField     = "vector"
	productVectorContextSparseField = "text_sparse_vector"

	productIdxProductVector           = "idx_product_vector"
	productIdxProductTextSparseVector = "idx_product_text_sparse_vector"

	// 以下长度限制用于控制写入 Milvus 的 varchar 字段大小，避免超出 schema 限制。
	productVectorNameMaxLength     = 512
	productVectorBarCodeMaxLength  = 128
	productVectorCategoryMaxLength = 256
	productVectorUnitMaxLength     = 128
	productVectorStandardMaxLength = 512
	productVectorRemarkMaxLength   = 2048
	productVectorContentMaxLength  = 16384

	// 以下检索参数用于控制初始召回规模与对外返回上限，给应用层重排预留足够候选集。
	productVectorSearchCandidateMultiplier = 3
	productVectorSearchMinCandidates       = 20
	productVectorSearchMaxCandidates       = 100
	productVectorSearchDefaultLimit        = 10
	productVectorSearchMaxLimit            = 50
)

type ProductVectorDaoImpl struct{}

type productVectorConfig struct {
	Collection string `mapstructure:"collection"`
}

// NewProductVectorDao 创建产品向量 DAO。
func NewProductVectorDao() masterdao.IProductVectorDao {
	return &ProductVectorDaoImpl{}
}
