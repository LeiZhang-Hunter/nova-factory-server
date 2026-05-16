package stockdao

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IStockMoveDao ERP 库存调拨单数据访问接口
type IStockMoveDao interface {
	Create(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error)
	Update(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockMove, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockMove, error)
	ListPage(c *gin.Context, req *stockmodels.StockMoveQuery) (*erpbiz.PageResult[stockmodels.StockMove], error)
	List(c *gin.Context, req *stockmodels.StockMoveQuery) (*stockmodels.StockMoveListData, error)
}
