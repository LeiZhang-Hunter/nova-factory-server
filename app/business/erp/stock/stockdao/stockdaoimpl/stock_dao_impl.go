package stockdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockDaoImpl 提供 ERP 产品库存数据访问能力。
type StockDaoImpl struct {
	*erpcrud.CRUDDao[stockmodels.Stock, stockmodels.StockUpsert, stockmodels.StockQuery]
}

// NewStockDao 创建 ERP 产品库存 DAO。
func NewStockDao(db *gorm.DB) stockdao.IStockDao {
	return &StockDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[stockmodels.Stock, stockmodels.StockUpsert, stockmodels.StockQuery](db, erpcrud.EntityConfig{
			Table:        "erp_stock",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 产品库存列表。
func (d *StockDaoImpl) List(c *gin.Context, req *stockmodels.StockQuery) (*stockmodels.StockListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockListData{Rows: result.Rows, Total: result.Total}, nil
}
