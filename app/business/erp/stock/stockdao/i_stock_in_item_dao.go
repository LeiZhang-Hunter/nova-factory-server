package stockdao

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IStockInItemDao ERP 其它入库单项数据访问接口
type IStockInItemDao interface {
	Create(c *gin.Context, req *stockmodels.StockInItemUpsert) (*stockmodels.StockInItem, error)
	Update(c *gin.Context, req *stockmodels.StockInItemUpsert) (*stockmodels.StockInItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockInItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockInItem, error)
	ListPage(c *gin.Context, req *stockmodels.StockInItemQuery) (*erpbiz.PageResult[stockmodels.StockInItem], error)
	List(c *gin.Context, req *stockmodels.StockInItemQuery) (*stockmodels.StockInItemListData, error)
}
