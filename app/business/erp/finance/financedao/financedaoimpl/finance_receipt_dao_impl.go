package financedaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FinanceReceiptDaoImpl 提供 ERP 收款单数据访问能力。
type FinanceReceiptDaoImpl struct {
	*erpcrud.CRUDDao[financemodels.FinanceReceipt, financemodels.FinanceReceiptUpsert, financemodels.FinanceReceiptQuery]
}

// NewFinanceReceiptDao 创建 ERP 收款单 DAO。
func NewFinanceReceiptDao(db *gorm.DB) financedao.IFinanceReceiptDao {
	return &FinanceReceiptDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[financemodels.FinanceReceipt, financemodels.FinanceReceiptUpsert, financemodels.FinanceReceiptQuery](db, erpcrud.EntityConfig{
			Table:        "erp_finance_receipt",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "收款单号"}},
		}),
	}
}

// List 查询 ERP 收款单列表。
func (d *FinanceReceiptDaoImpl) List(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*financemodels.FinanceReceiptListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinanceReceiptListData{Rows: result.Rows, Total: result.Total}, nil
}
