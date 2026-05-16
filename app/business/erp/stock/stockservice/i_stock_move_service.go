package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockMoveService ERP 库存调拨单服务接口
type IStockMoveService interface {
	Create(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error)
	Update(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockMove, error)
	List(c *gin.Context, req *stockmodels.StockMoveQuery) (*stockmodels.StockMoveListData, error)
}
