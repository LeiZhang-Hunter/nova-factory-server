package purchaseserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"

	"github.com/gin-gonic/gin"
)

// PurchaseReturnItemServiceImpl 提供 ERP 采购退货项业务实现。
type PurchaseReturnItemServiceImpl struct {
	*erpcrud.CRUDService[purchasemodels.PurchaseReturnItem, purchasemodels.PurchaseReturnItemUpsert, purchasemodels.PurchaseReturnItemQuery]
}

// NewPurchaseReturnItemService 创建 ERP 采购退货项服务。
func NewPurchaseReturnItemService(dao purchasedao.IPurchaseReturnItemDao) purchaseservice.IPurchaseReturnItemService {
	return &PurchaseReturnItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[purchasemodels.PurchaseReturnItem, purchasemodels.PurchaseReturnItemUpsert, purchasemodels.PurchaseReturnItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_purchase_return_items",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 采购退货项列表。
func (s *PurchaseReturnItemServiceImpl) List(c *gin.Context, req *purchasemodels.PurchaseReturnItemQuery) (*purchasemodels.PurchaseReturnItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseReturnItemListData{Rows: result.Rows, Total: result.Total}, nil
}
