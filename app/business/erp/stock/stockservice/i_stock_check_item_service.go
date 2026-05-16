package stockservice

import (
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
)

// IStockCheckItemService ERP 库存盘点单项服务接口
type IStockCheckItemService interface {
	Create(c *gin.Context, req *stockmodels.StockCheckItemUpsert) (*stockmodels.StockCheckItem, error)
	Update(c *gin.Context, req *stockmodels.StockCheckItemUpsert) (*stockmodels.StockCheckItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*stockmodels.StockCheckItem, error)
	List(c *gin.Context, req *stockmodels.StockCheckItemQuery) (*stockmodels.StockCheckItemListData, error)
}
