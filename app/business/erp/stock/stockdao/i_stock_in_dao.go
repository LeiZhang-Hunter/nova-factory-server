package stockdao

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockInDao ERP 其它入库单数据访问接口
type IStockInDao interface {
	Create(c *gin.Context, req *stockmodels.StockInUpsert) (*stockmodels.StockIn, error)
	Update(c *gin.Context, req *stockmodels.StockInUpsert) (*stockmodels.StockIn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockIn, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockIn, error)
	ListPage(c *gin.Context, req *stockmodels.StockInQuery) (*stockmodels.StockInListData, error)
	List(c *gin.Context, req *stockmodels.StockInQuery) (*stockmodels.StockInListData, error)
}
