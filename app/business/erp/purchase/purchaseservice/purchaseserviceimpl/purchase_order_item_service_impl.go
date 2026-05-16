package purchaseserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"

	"github.com/gin-gonic/gin"
)

// PurchaseOrderItemServiceImpl 提供 ERP 采购订单项业务实现。
type PurchaseOrderItemServiceImpl struct {
	*erpcrud.CRUDService[purchasemodels.PurchaseOrderItem, purchasemodels.PurchaseOrderItemUpsert, purchasemodels.PurchaseOrderItemQuery]
}

// NewPurchaseOrderItemService 创建 ERP 采购订单项服务。
func NewPurchaseOrderItemService(dao purchasedao.IPurchaseOrderItemDao) purchaseservice.IPurchaseOrderItemService {
	return &PurchaseOrderItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[purchasemodels.PurchaseOrderItem, purchasemodels.PurchaseOrderItemUpsert, purchasemodels.PurchaseOrderItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_purchase_order_items",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 采购订单项列表。
func (s *PurchaseOrderItemServiceImpl) List(c *gin.Context, req *purchasemodels.PurchaseOrderItemQuery) (*purchasemodels.PurchaseOrderItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseOrderItemListData{Rows: result.Rows, Total: result.Total}, nil
}
