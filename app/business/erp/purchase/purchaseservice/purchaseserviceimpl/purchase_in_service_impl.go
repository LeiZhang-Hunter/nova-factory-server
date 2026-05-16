package purchaseserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"

	"github.com/gin-gonic/gin"
)

// PurchaseInServiceImpl 提供 ERP 采购入库业务实现。
type PurchaseInServiceImpl struct {
	*erpcrud.CRUDService[purchasemodels.PurchaseIn, purchasemodels.PurchaseInUpsert, purchasemodels.PurchaseInQuery]
}

// NewPurchaseInService 创建 ERP 采购入库服务。
func NewPurchaseInService(dao purchasedao.IPurchaseInDao) purchaseservice.IPurchaseInService {
	return &PurchaseInServiceImpl{
		CRUDService: erpcrud.NewCRUDService[purchasemodels.PurchaseIn, purchasemodels.PurchaseInUpsert, purchasemodels.PurchaseInQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_purchase_in",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "采购入库单号"}},
		}),
	}
}

// List 查询 ERP 采购入库列表。
func (s *PurchaseInServiceImpl) List(c *gin.Context, req *purchasemodels.PurchaseInQuery) (*purchasemodels.PurchaseInListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseInListData{Rows: result.Rows, Total: result.Total}, nil
}
