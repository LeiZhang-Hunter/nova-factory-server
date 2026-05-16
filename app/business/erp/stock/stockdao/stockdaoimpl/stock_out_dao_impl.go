package stockdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockOutDaoImpl 提供 ERP 其它出库单数据访问能力。
type StockOutDaoImpl struct {
	*erpcrud.CRUDDao[stockmodels.StockOut, stockmodels.StockOutUpsert, stockmodels.StockOutQuery]
}

// NewStockOutDao 创建 ERP 其它出库单 DAO。
func NewStockOutDao(db *gorm.DB) stockdao.IStockOutDao {
	return &StockOutDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[stockmodels.StockOut, stockmodels.StockOutUpsert, stockmodels.StockOutQuery](db, erpcrud.EntityConfig{
			Table:        "erp_stock_out",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "出库单号"}},
		}),
	}
}

// List 查询 ERP 其它出库单列表。
func (d *StockOutDaoImpl) List(c *gin.Context, req *stockmodels.StockOutQuery) (*stockmodels.StockOutListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockOutListData{Rows: result.Rows, Total: result.Total}, nil
}
