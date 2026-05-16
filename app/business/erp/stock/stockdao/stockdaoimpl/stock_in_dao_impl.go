package stockdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockInDaoImpl 提供 ERP 其它入库单数据访问能力。
type StockInDaoImpl struct {
	*erpcrud.CRUDDao[stockmodels.StockIn, stockmodels.StockInUpsert, stockmodels.StockInQuery]
}

// NewStockInDao 创建 ERP 其它入库单 DAO。
func NewStockInDao(db *gorm.DB) stockdao.IStockInDao {
	return &StockInDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[stockmodels.StockIn, stockmodels.StockInUpsert, stockmodels.StockInQuery](db, erpcrud.EntityConfig{
			Table:        "erp_stock_in",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "入库单号"}},
		}),
	}
}

// List 查询 ERP 其它入库单列表。
func (d *StockInDaoImpl) List(c *gin.Context, req *stockmodels.StockInQuery) (*stockmodels.StockInListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockInListData{Rows: result.Rows, Total: result.Total}, nil
}
