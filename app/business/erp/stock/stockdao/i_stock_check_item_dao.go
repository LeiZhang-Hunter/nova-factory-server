package stockdao

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IStockCheckItemDao ERP 库存盘点单项数据访问接口
type IStockCheckItemDao interface {
	Create(c *gin.Context, req *stockmodels.StockCheckItemUpsert) (*stockmodels.StockCheckItem, error)
	Update(c *gin.Context, req *stockmodels.StockCheckItemUpsert) (*stockmodels.StockCheckItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockCheckItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*stockmodels.StockCheckItem, error)
	ListPage(c *gin.Context, req *stockmodels.StockCheckItemQuery) (*erpbiz.PageResult[stockmodels.StockCheckItem], error)
	List(c *gin.Context, req *stockmodels.StockCheckItemQuery) (*stockmodels.StockCheckItemListData, error)
}
