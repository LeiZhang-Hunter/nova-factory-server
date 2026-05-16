package financeserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"
	"nova-factory-server/app/business/erp/finance/financeservice"

	"github.com/gin-gonic/gin"
)

// FinancePaymentItemServiceImpl 提供 ERP 付款项业务实现。
type FinancePaymentItemServiceImpl struct {
	*erpcrud.CRUDService[financemodels.FinancePaymentItem, financemodels.FinancePaymentItemUpsert, financemodels.FinancePaymentItemQuery]
}

// NewFinancePaymentItemService 创建 ERP 付款项服务。
func NewFinancePaymentItemService(dao financedao.IFinancePaymentItemDao) financeservice.IFinancePaymentItemService {
	return &FinancePaymentItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[financemodels.FinancePaymentItem, financemodels.FinancePaymentItemUpsert, financemodels.FinancePaymentItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_finance_payment_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 付款项列表。
func (s *FinancePaymentItemServiceImpl) List(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*financemodels.FinancePaymentItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinancePaymentItemListData{Rows: result.Rows, Total: result.Total}, nil
}
