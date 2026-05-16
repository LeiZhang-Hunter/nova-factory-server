package stockdao

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockOutItemDao ERP 其它出库单项数据访问接口
type IStockOutItemDao interface {
	Create(c *gin.Context, req *stockmodels.StockOutItemUpsert) (*stockmodels.StockOutItem, error)
	Update(c *gin.Context, req *stockmodels.StockOutItemUpsert) (*stockmodels.StockOutItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockOutItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockOutItem, error)
	ListPage(c *gin.Context, req *stockmodels.StockOutItemQuery) (*stockmodels.StockOutItemListData, error)
	List(c *gin.Context, req *stockmodels.StockOutItemQuery) (*stockmodels.StockOutItemListData, error)
}
