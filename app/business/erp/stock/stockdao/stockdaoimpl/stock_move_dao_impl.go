package stockdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockMoveDaoImpl 提供 ERP 库存调拨单数据访问能力。
type StockMoveDaoImpl struct {
	*erpcrud.CRUDDao[stockmodels.StockMove, stockmodels.StockMoveUpsert, stockmodels.StockMoveQuery]
}

// NewStockMoveDao 创建 ERP 库存调拨单 DAO。
func NewStockMoveDao(db *gorm.DB) stockdao.IStockMoveDao {
	return &StockMoveDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[stockmodels.StockMove, stockmodels.StockMoveUpsert, stockmodels.StockMoveQuery](db, erpcrud.EntityConfig{
			Table:        "erp_stock_move",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "调拨单号"}},
		}),
	}
}

// List 查询 ERP 库存调拨单列表。
func (d *StockMoveDaoImpl) List(c *gin.Context, req *stockmodels.StockMoveQuery) (*stockmodels.StockMoveListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockMoveListData{Rows: result.Rows, Total: result.Total}, nil
}
