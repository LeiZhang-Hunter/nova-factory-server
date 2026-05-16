package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockService ERP 产品库存服务接口
type IStockService interface {
	Create(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error)
	Update(c *gin.Context, req *stockmodels.StockUpsert) (*stockmodels.Stock, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.Stock, error)
	List(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error)
}
