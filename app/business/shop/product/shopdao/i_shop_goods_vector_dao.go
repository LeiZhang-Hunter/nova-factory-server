package shopdao

import (
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

// IShopGoodsVectorDao 商品向量数据访问接口
type IShopGoodsVectorDao interface {
	// Upsert 将商品及其 SKU 向量写入 Milvus。
	// 每条 SKU 会落为一行，便于检索时保留更细的规格粒度。
	Upsert(c *gin.Context, goods *shopmodels.Goods, items []*shopmodels.GoodsVectorUpsertItem) (*shopmodels.GoodsVectorResult, error)

	// Search 复用批量检索入口处理单条查询。
	Search(c *gin.Context, req *shopmodels.GoodsVectorSearchReq, vector []float32, fallbackWithoutMetadata bool) (*shopmodels.GoodsVectorSearchData, error)

	// BatchSearch 执行商品批量向量检索。
	// 支持 dense-only 检索，以及 dense + BM25 稀疏向量的 hybrid 检索。
	BatchSearch(c *gin.Context, req *shopmodels.GoodsVectorBatchSearchReq, vectors [][]float32, fallbackWithoutMetadata bool) (*shopmodels.GoodsVectorBatchSearchData, error)
}
