package stockdao

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockCheckDao ERP 库存盘点单数据访问接口
type IStockCheckDao interface {
	Create(c *gin.Context, req *stockmodels.StockCheckUpsert) (*stockmodels.StockCheck, error)
	Update(c *gin.Context, req *stockmodels.StockCheckUpsert) (*stockmodels.StockCheck, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockCheck, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockCheck, error)
	ListPage(c *gin.Context, req *stockmodels.StockCheckQuery) (*stockmodels.StockCheckListData, error)
	List(c *gin.Context, req *stockmodels.StockCheckQuery) (*stockmodels.StockCheckListData, error)
}
