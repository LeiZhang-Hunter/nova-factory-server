package purchaseserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"

	"github.com/gin-gonic/gin"
)

// PurchaseInItemServiceImpl 提供 ERP 采购入库项业务实现。
type PurchaseInItemServiceImpl struct {
	*erpcrud.CRUDService[purchasemodels.PurchaseInItem, purchasemodels.PurchaseInItemUpsert, purchasemodels.PurchaseInItemQuery]
}

// NewPurchaseInItemService 创建 ERP 采购入库项服务。
func NewPurchaseInItemService(dao purchasedao.IPurchaseInItemDao) purchaseservice.IPurchaseInItemService {
	return &PurchaseInItemServiceImpl{
		CRUDService: erpcrud.NewCRUDService[purchasemodels.PurchaseInItem, purchasemodels.PurchaseInItemUpsert, purchasemodels.PurchaseInItemQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_purchase_in_items",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 采购入库项列表。
func (s *PurchaseInItemServiceImpl) List(c *gin.Context, req *purchasemodels.PurchaseInItemQuery) (*purchasemodels.PurchaseInItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseInItemListData{Rows: result.Rows, Total: result.Total}, nil
}
