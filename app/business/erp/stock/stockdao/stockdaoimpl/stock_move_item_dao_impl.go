package stockdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockMoveItemDaoImpl 提供 ERP 库存调拨单项数据访问能力。
type StockMoveItemDaoImpl struct {
	*erpcrud.CRUDDao[stockmodels.StockMoveItem, stockmodels.StockMoveItemUpsert, stockmodels.StockMoveItemQuery]
}

// NewStockMoveItemDao 创建 ERP 库存调拨单项 DAO。
func NewStockMoveItemDao(db *gorm.DB) stockdao.IStockMoveItemDao {
	return &StockMoveItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[stockmodels.StockMoveItem, stockmodels.StockMoveItemUpsert, stockmodels.StockMoveItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_stock_move_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 库存调拨单项列表。
func (d *StockMoveItemDaoImpl) List(c *gin.Context, req *stockmodels.StockMoveItemQuery) (*stockmodels.StockMoveItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockMoveItemListData{Rows: result.Rows, Total: result.Total}, nil
}
