package shopservice

import (
	"encoding/csv"
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

// IShopGoodsService 商品服务接口
type IShopGoodsService interface {
	Create(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error)
	Update(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error)
	GetByGoodsID(c *gin.Context, goodsID string) (*shopmodels.Goods, error)
	List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error)
	ExportCSV(c *gin.Context, req *shopmodels.GoodsQuery, csvWriter *csv.Writer, flush func()) error
	Import(c *gin.Context, records []shopmodels.ImportGoodsRecord) error
	GenerateVector(c *gin.Context, req *shopmodels.GenVectorReq) (*shopmodels.GoodsVectorResult, error)
	GenerateAllOnSaleVectors(c *gin.Context, req *shopmodels.GenAllVectorReq) (*shopmodels.GoodsVectorTaskData, error)
	GetGenerateAllOnSaleVectorsProgress(c *gin.Context, taskID string) (*shopmodels.GoodsVectorTaskProgress, error)
	ListGenerateAllOnSaleVectorTasks(c *gin.Context) (*shopmodels.GoodsVectorTaskListData, error)
	SearchVector(c *gin.Context, req *shopmodels.GoodsVectorSearchReq) (*shopmodels.GoodsVectorSearchData, error)
	BatchSearchVector(c *gin.Context, req *shopmodels.GoodsVectorBatchSearchReq) (*shopmodels.GoodsVectorBatchSearchData, error)
}
