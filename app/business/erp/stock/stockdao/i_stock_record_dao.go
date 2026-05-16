package stockdao

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockRecordDao ERP 产品库存明细数据访问接口
type IStockRecordDao interface {
	Create(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error)
	Update(c *gin.Context, req *stockmodels.StockRecordUpsert) (*stockmodels.StockRecord, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockRecord, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockRecord, error)
	ListPage(c *gin.Context, req *stockmodels.StockRecordQuery) (*stockmodels.StockRecordListData, error)
	List(c *gin.Context, req *stockmodels.StockRecordQuery) (*stockmodels.StockRecordListData, error)
}
