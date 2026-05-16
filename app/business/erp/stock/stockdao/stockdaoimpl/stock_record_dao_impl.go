package stockdaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StockRecordDaoImpl 提供 ERP 产品库存明细数据访问能力。
type StockRecordDaoImpl struct {
	*erpcrud.CRUDDao[stockmodels.StockRecord, stockmodels.StockRecordUpsert, stockmodels.StockRecordQuery]
}

// NewStockRecordDao 创建 ERP 产品库存明细 DAO。
func NewStockRecordDao(db *gorm.DB) stockdao.IStockRecordDao {
	return &StockRecordDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[stockmodels.StockRecord, stockmodels.StockRecordUpsert, stockmodels.StockRecordQuery](db, erpcrud.EntityConfig{
			Table:        "erp_stock_record",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 产品库存明细列表。
func (d *StockRecordDaoImpl) List(c *gin.Context, req *stockmodels.StockRecordQuery) (*stockmodels.StockRecordListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockRecordListData{Rows: result.Rows, Total: result.Total}, nil
}
