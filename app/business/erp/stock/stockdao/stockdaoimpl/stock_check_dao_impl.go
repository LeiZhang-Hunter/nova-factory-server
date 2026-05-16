package stockdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockCheckDaoImpl 提供 ERP 库存盘点单数据访问能力。
type StockCheckDaoImpl struct {
	*erpcrud.CRUDDao[stockmodels.StockCheck, stockmodels.StockCheckUpsert, stockmodels.StockCheckQuery]
}

// NewStockCheckDao 创建 ERP 库存盘点单 DAO。
func NewStockCheckDao(db *gorm.DB) stockdao.IStockCheckDao {
	return &StockCheckDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[stockmodels.StockCheck, stockmodels.StockCheckUpsert, stockmodels.StockCheckQuery](db, erpcrud.EntityConfig{
			Table:        "erp_stock_check",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "盘点单号"}},
		}),
	}
}

// List 查询 ERP 库存盘点单列表。
func (d *StockCheckDaoImpl) List(c *gin.Context, req *stockmodels.StockCheckQuery) (*stockmodels.StockCheckListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockCheckListData{Rows: result.Rows, Total: result.Total}, nil
}
