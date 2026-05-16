package stockdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockOutItemDaoImpl 提供 ERP 其它出库单项数据访问能力。
type StockOutItemDaoImpl struct {
	*erpcrud.CRUDDao[stockmodels.StockOutItem, stockmodels.StockOutItemUpsert, stockmodels.StockOutItemQuery]
}

// NewStockOutItemDao 创建 ERP 其它出库单项 DAO。
func NewStockOutItemDao(db *gorm.DB) stockdao.IStockOutItemDao {
	return &StockOutItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[stockmodels.StockOutItem, stockmodels.StockOutItemUpsert, stockmodels.StockOutItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_stock_out_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 其它出库单项列表。
func (d *StockOutItemDaoImpl) List(c *gin.Context, req *stockmodels.StockOutItemQuery) (*stockmodels.StockOutItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockOutItemListData{Rows: result.Rows, Total: result.Total}, nil
}
