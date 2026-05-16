package financeserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"
	"nova-factory-server/app/business/erp/finance/financeservice"

	"github.com/gin-gonic/gin"
)

// FinanceReceiptItemServiceImpl 提供 ERP 收款项业务实现。
type FinanceReceiptItemServiceImpl struct {
	*erpcrud.CRUDService[financemodels.FinanceReceiptItem, financemodels.FinanceReceiptItemUpsert, financemodels.FinanceReceiptItemQuery]
}

// NewFinanceReceiptItemService 创建 ERP 收款项服务。
func NewFinanceReceiptItemService(dao financedao.IFinanceReceiptItemDao) financeservice.IFinanceReceiptItemService {
	return &FinanceReceiptItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[financemodels.FinanceReceiptItem, financemodels.FinanceReceiptItemUpsert, financemodels.FinanceReceiptItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_finance_receipt_item",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 收款项列表。
func (s *FinanceReceiptItemServiceImpl) List(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*financemodels.FinanceReceiptItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinanceReceiptItemListData{Rows: result.Rows, Total: result.Total}, nil
}
