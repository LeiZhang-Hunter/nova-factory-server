package stockdao

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IStockDao ERP 产品库存数据访问接口
type IStockDao interface {
	Create(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error)
	Update(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.Stock, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.Stock, error)
	ListPage(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error)
	List(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error)
	UpdateStockByProductIDWithDB(db *gorm.DB, productID int64, count float64) error
	UpdateStockByID(c *gin.Context, id int64, count float64) error
	UpsertByID(c *gin.Context, id int64, updates map[string]any) error
}
