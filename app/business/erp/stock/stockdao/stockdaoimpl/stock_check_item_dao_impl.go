package stockdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockCheckItemDaoImpl 提供 ERP 库存盘点单项数据访问能力。
type StockCheckItemDaoImpl struct {
	*erpcrud.CRUDDao[stockmodels.StockCheckItem, stockmodels.StockCheckItemUpsert, stockmodels.StockCheckItemQuery]
}

// NewStockCheckItemDao 创建 ERP 库存盘点单项 DAO。
func NewStockCheckItemDao(db *gorm.DB) stockdao.IStockCheckItemDao {
	return &StockCheckItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[stockmodels.StockCheckItem, stockmodels.StockCheckItemUpsert, stockmodels.StockCheckItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_stock_check_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 库存盘点单项列表。
func (d *StockCheckItemDaoImpl) List(c *gin.Context, req *stockmodels.StockCheckItemQuery) (*stockmodels.StockCheckItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockCheckItemListData{Rows: result.Rows, Total: result.Total}, nil
}
