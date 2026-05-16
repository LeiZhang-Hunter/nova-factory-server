package stockdao

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockOutDao ERP 其它出库单数据访问接口
type IStockOutDao interface {
	Create(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error)
	Update(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockOut, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockOut, error)
	ListPage(c *gin.Context, req *stockmodels.StockOutQuery) (*stockmodels.StockOutListData, error)
	List(c *gin.Context, req *stockmodels.StockOutQuery) (*stockmodels.StockOutListData, error)
}
