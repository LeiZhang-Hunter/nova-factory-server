package financedaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FinanceReceiptItemDaoImpl 提供 ERP 收款项数据访问能力。
type FinanceReceiptItemDaoImpl struct {
	*erpcrud.CRUDDao[financemodels.FinanceReceiptItem, financemodels.FinanceReceiptItemUpsert, financemodels.FinanceReceiptItemQuery]
}

// NewFinanceReceiptItemDao 创建 ERP 收款项 DAO。
func NewFinanceReceiptItemDao(db *gorm.DB) financedao.IFinanceReceiptItemDao {
	return &FinanceReceiptItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[financemodels.FinanceReceiptItem, financemodels.FinanceReceiptItemUpsert, financemodels.FinanceReceiptItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_finance_receipt_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 收款项列表。
func (d *FinanceReceiptItemDaoImpl) List(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*financemodels.FinanceReceiptItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinanceReceiptItemListData{Rows: result.Rows, Total: result.Total}, nil
}
