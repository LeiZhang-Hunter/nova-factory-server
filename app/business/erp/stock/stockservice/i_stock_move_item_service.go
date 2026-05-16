package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockMoveItemService ERP 库存调拨单项服务接口
type IStockMoveItemService interface {
	Create(c *gin.Context, req *stockmodels.StockMoveItemUpsert) (*stockmodels.StockMoveItem, error)
	Update(c *gin.Context, req *stockmodels.StockMoveItemUpsert) (*stockmodels.StockMoveItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockMoveItem, error)
	List(c *gin.Context, req *stockmodels.StockMoveItemQuery) (*stockmodels.StockMoveItemListData, error)
}
