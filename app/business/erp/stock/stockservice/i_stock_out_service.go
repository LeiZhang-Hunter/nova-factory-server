package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockOutService ERP 其它出库单服务接口
type IStockOutService interface {
	Create(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error)
	Update(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockOut, error)
	List(c *gin.Context, req *stockmodels.StockOutQuery) (*stockmodels.StockOutListData, error)
}
