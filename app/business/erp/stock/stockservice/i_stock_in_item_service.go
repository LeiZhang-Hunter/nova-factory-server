package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockInItemService ERP 其它入库单项服务接口
type IStockInItemService interface {
	Create(c *gin.Context, req *stockmodels.StockInItemUpsert) (*stockmodels.StockInItem, error)
	Update(c *gin.Context, req *stockmodels.StockInItemUpsert) (*stockmodels.StockInItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockInItem, error)
	List(c *gin.Context, req *stockmodels.StockInItemQuery) (*stockmodels.StockInItemListData, error)
}
