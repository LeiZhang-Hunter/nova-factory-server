package financedaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FinancePaymentItemDaoImpl 提供 ERP 付款项数据访问能力。
type FinancePaymentItemDaoImpl struct {
	*erpcrud.CRUDDao[financemodels.FinancePaymentItem, financemodels.FinancePaymentItemUpsert, financemodels.FinancePaymentItemQuery]
}

// NewFinancePaymentItemDao 创建 ERP 付款项 DAO。
func NewFinancePaymentItemDao(db *gorm.DB) financedao.IFinancePaymentItemDao {
	return &FinancePaymentItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[financemodels.FinancePaymentItem, financemodels.FinancePaymentItemUpsert, financemodels.FinancePaymentItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_finance_payment_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 付款项列表。
func (d *FinancePaymentItemDaoImpl) List(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*financemodels.FinancePaymentItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinancePaymentItemListData{Rows: result.Rows, Total: result.Total}, nil
}
