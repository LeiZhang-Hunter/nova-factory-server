package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockCheckService ERP 库存盘点单服务接口
type IStockCheckService interface {
	Create(c *gin.Context, req *stockmodels.StockCheckUpsert) (*stockmodels.StockCheck, error)
	Update(c *gin.Context, req *stockmodels.StockCheckUpsert) (*stockmodels.StockCheck, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockCheck, error)
	List(c *gin.Context, req *stockmodels.StockCheckQuery) (*stockmodels.StockCheckListData, error)
}
