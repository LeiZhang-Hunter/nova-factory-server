package financedaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FinancePaymentDaoImpl 提供 ERP 付款单数据访问能力。
type FinancePaymentDaoImpl struct {
	*erpcrud.CRUDDao[financemodels.FinancePayment, financemodels.FinancePaymentUpsert, financemodels.FinancePaymentQuery]
}

// NewFinancePaymentDao 创建 ERP 付款单 DAO。
func NewFinancePaymentDao(db *gorm.DB) financedao.IFinancePaymentDao {
	return &FinancePaymentDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[financemodels.FinancePayment, financemodels.FinancePaymentUpsert, financemodels.FinancePaymentQuery](db, erpcrud.EntityConfig{
			Table:        "erp_finance_payment",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "付款单号"}},
		}),
	}
}

// List 查询 ERP 付款单列表。
func (d *FinancePaymentDaoImpl) List(c *gin.Context, req *financemodels.FinancePaymentQuery) (*financemodels.FinancePaymentListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinancePaymentListData{Rows: result.Rows, Total: result.Total}, nil
}
