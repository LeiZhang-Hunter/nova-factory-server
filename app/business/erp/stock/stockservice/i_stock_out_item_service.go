package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockOutItemService ERP 其它出库单项服务接口
type IStockOutItemService interface {
	Create(c *gin.Context, req *stockmodels.StockOutItemUpsert) (*stockmodels.StockOutItem, error)
	Update(c *gin.Context, req *stockmodels.StockOutItemUpsert) (*stockmodels.StockOutItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockOutItem, error)
	List(c *gin.Context, req *stockmodels.StockOutItemQuery) (*stockmodels.StockOutItemListData, error)
}
