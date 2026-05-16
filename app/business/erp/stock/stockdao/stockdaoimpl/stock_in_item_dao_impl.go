package stockdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockInItemDaoImpl 提供 ERP 其它入库单项数据访问能力。
type StockInItemDaoImpl struct {
	*erpcrud.CRUDDao[stockmodels.StockInItem, stockmodels.StockInItemUpsert, stockmodels.StockInItemQuery]
}

// NewStockInItemDao 创建 ERP 其它入库单项 DAO。
func NewStockInItemDao(db *gorm.DB) stockdao.IStockInItemDao {
	return &StockInItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[stockmodels.StockInItem, stockmodels.StockInItemUpsert, stockmodels.StockInItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_stock_in_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 其它入库单项列表。
func (d *StockInItemDaoImpl) List(c *gin.Context, req *stockmodels.StockInItemQuery) (*stockmodels.StockInItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockInItemListData{Rows: result.Rows, Total: result.Total}, nil
}
