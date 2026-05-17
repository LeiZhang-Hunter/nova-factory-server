package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockInService ERP 其它入库单服务接口
type IStockInService interface {
	Create(c *gin.Context, req *stockmodels.StockInUpsert) (*stockmodels.StockIn, error)
	Update(c *gin.Context, req *stockmodels.StockInUpsert) (*stockmodels.StockIn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockIn, error)
	List(c *gin.Context, req *stockmodels.StockInQuery) (*stockmodels.StockInListData, error)
}
