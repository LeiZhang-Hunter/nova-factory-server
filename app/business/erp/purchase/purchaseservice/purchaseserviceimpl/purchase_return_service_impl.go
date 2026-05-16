package purchaseserviceimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"

	"github.com/gin-gonic/gin"
)

// PurchaseReturnServiceImpl 提供 ERP 采购退货业务实现。
type PurchaseReturnServiceImpl struct {
	*erpcrud.CRUDService[purchasemodels.PurchaseReturn, purchasemodels.PurchaseReturnUpsert, purchasemodels.PurchaseReturnQuery]
}

// NewPurchaseReturnService 创建 ERP 采购退货服务。
func NewPurchaseReturnService(dao purchasedao.IPurchaseReturnDao) purchaseservice.IPurchaseReturnService {
	return &PurchaseReturnServiceImpl{
		CRUDService: erpcrud.NewCRUDService[purchasemodels.PurchaseReturn, purchasemodels.PurchaseReturnUpsert, purchasemodels.PurchaseReturnQuery](dao, erpcrud.EntityConfig{
			Table:        "erp_purchase_return",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "采购退货单号"}},
		}),
	}
}

// List 查询 ERP 采购退货列表。
func (s *PurchaseReturnServiceImpl) List(c *gin.Context, req *purchasemodels.PurchaseReturnQuery) (*purchasemodels.PurchaseReturnListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseReturnListData{Rows: result.Rows, Total: result.Total}, nil
}
