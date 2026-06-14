package shopservice

import (
	"encoding/csv"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"

	"github.com/gin-gonic/gin"
)

// IShopGoodsService 定义商品领域的应用服务能力。
// 该接口覆盖商品基础 CRUD、导入导出、向量生成、任务进度查询以及向量检索等能力，
// 供 API 层按统一入口调用。
type IShopGoodsService interface {
	// Create 创建商品基础信息。
	Create(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error)

	// Update 更新商品基础信息。
	Update(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error)

	// DeleteByIDs 按数据库主键批量删除商品。
	DeleteByIDs(c *gin.Context, ids []int64) error

	// GetByID 按数据库主键查询单个商品详情。
	GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error)

	// GetByGoodsID 按业务商品 ID 查询单个商品详情。
	GetByGoodsID(c *gin.Context, goodsID string) (*shopmodels.Goods, error)

	// List 分页查询商品列表。
	List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error)

	// ExportCSV 按查询条件导出商品 CSV。
	// csvWriter 由调用方传入，flush 用于在流式写出时主动刷新缓冲。
	ExportCSV(c *gin.Context, req *shopmodels.GoodsQuery, csvWriter *csv.Writer, flush func()) error

	// Import 批量导入商品及其 SKU 数据。
	Import(c *gin.Context, records []shopmodels.ImportGoodsRecord) error

	// GenerateVector 为单个商品生成并写入向量数据。
	GenerateVector(c *gin.Context, req *shopmodels.GenVectorReq) (*shopmodels.GoodsVectorResult, error)

	// GenerateAllOnSaleVectors 为全部在售商品触发异步向量生成任务。
	GenerateAllOnSaleVectors(c *gin.Context, req *shopmodels.GenAllVectorReq) (*shopmodels.GoodsVectorTaskData, error)

	// GetGenerateAllOnSaleVectorsProgress 查询全量在售商品向量生成任务的进度。
	GetGenerateAllOnSaleVectorsProgress(c *gin.Context, taskID string) (*shopmodels.GoodsVectorTaskProgress, error)

	// ListGenerateAllOnSaleVectorTasks 列出全量在售商品向量生成任务记录。
	ListGenerateAllOnSaleVectorTasks(c *gin.Context) (*shopmodels.GoodsVectorTaskListData, error)

	// SearchVector 执行单条商品向量检索。
	SearchVector(c *gin.Context, req *shopmodels.GoodsVectorSearchReq) (*shopmodels.GoodsVectorSearchData, error)

	// BatchSearchVector 执行批量商品向量检索。
	BatchSearchVector(c *gin.Context, req *shopmodels.GoodsVectorBatchSearchReq) (*shopmodels.GoodsVectorBatchSearchData, error)

	// SyncEvent 同步事件
	SyncEvent(event event.ProductEvent) (result.SyncProductResponse, error)

	// SyncStock 同步库存变更，使用传入的 db 保证事务一致性
	SyncStock(stocks event.StockEvent) error
}
