package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockRecordService ERP 产品库存明细服务接口
type IStockRecordService interface {
	Create(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error)
	Update(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockRecord, error)
	List(c *gin.Context, req *stockmodels.StockRecordQuery) (*stockmodels.StockRecordListData, error)
}
