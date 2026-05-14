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
	List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error)
	ExportCSV(c *gin.Context, req *shopmodels.GoodsQuery, csvWriter *csv.Writer, flush func()) error
	Import(c *gin.Context, records []shopmodels.ImportGoodsRecord) error
	GenerateVector(c *gin.Context, id int64) (*shopmodels.GoodsVectorResult, error)
}
