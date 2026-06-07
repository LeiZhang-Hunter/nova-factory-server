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

	// UpdateSaleStatusByGoodsID 按商品主键同步 Milvus 中所有向量行的在售状态。
	// 仅覆盖 is_sale 相关字段，不重新生成 embedding 内容。
	UpdateSaleStatusByGoodsID(c *gin.Context, goodsDBID int64, isOnSale int32) error

	// DeleteBySkuIDs 按 SKU 主键批量删除商品向量行。
	// 当前商品向量表按 SKU 落行，因此删除规格时应同步删除对应行。
	DeleteBySkuIDs(c *gin.Context, skuIDs []int64) error

	// Search 复用批量检索入口处理单条查询。
	Search(c *gin.Context, req *shopmodels.GoodsVectorSearchReq, vector []float32, fallbackWithoutMetadata bool) (*shopmodels.GoodsVectorSearchData, error)

	// BatchSearch 执行商品批量向量检索。
	// 支持 dense-only 检索，以及 dense + BM25 稀疏向量的 hybrid 检索。
	BatchSearch(c *gin.Context, req *shopmodels.GoodsVectorBatchSearchReq, vectors [][]float32, fallbackWithoutMetadata bool) (*shopmodels.GoodsVectorBatchSearchData, error)
}
