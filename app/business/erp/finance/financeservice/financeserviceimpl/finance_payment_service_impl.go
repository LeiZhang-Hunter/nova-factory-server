package financeserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"
	"nova-factory-server/app/business/erp/finance/financeservice"

	"github.com/gin-gonic/gin"
)

// FinancePaymentServiceImpl 提供 ERP 付款单业务实现。
type FinancePaymentServiceImpl struct {
	*erpcrud.CRUDService[financemodels.FinancePayment, financemodels.FinancePaymentUpsert, financemodels.FinancePaymentQuery]
}

// NewFinancePaymentService 创建 ERP 付款单服务。
func NewFinancePaymentService(dao financedao.IFinancePaymentDao) financeservice.IFinancePaymentService {
	return &FinancePaymentServiceImpl{
		CRUDService: erpcrud.NewCRUDService[financemodels.FinancePayment, financemodels.FinancePaymentUpsert, financemodels.FinancePaymentQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_finance_payment",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "付款单号"}},
		}),
	}
}

// List 查询 ERP 付款单列表。
func (s *FinancePaymentServiceImpl) List(c *gin.Context, req *financemodels.FinancePaymentQuery) (*financemodels.FinancePaymentListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinancePaymentListData{Rows: result.Rows, Total: result.Total}, nil
}
