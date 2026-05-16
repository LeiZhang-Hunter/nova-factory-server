package financeserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"
	"nova-factory-server/app/business/erp/finance/financeservice"

	"github.com/gin-gonic/gin"
)

// FinanceReceiptServiceImpl 提供 ERP 收款单业务实现。
type FinanceReceiptServiceImpl struct {
	*erpcrud.CRUDService[financemodels.FinanceReceipt, financemodels.FinanceReceiptUpsert, financemodels.FinanceReceiptQuery]
}

// NewFinanceReceiptService 创建 ERP 收款单服务。
func NewFinanceReceiptService(dao financedao.IFinanceReceiptDao) financeservice.IFinanceReceiptService {
	return &FinanceReceiptServiceImpl{
		CRUDService: erpcrud.NewCRUDService[financemodels.FinanceReceipt, financemodels.FinanceReceiptUpsert, financemodels.FinanceReceiptQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_finance_receipt",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "收款单号"}},
		}),
	}
}

// List 查询 ERP 收款单列表。
func (s *FinanceReceiptServiceImpl) List(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*financemodels.FinanceReceiptListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinanceReceiptListData{Rows: result.Rows, Total: result.Total}, nil
}
