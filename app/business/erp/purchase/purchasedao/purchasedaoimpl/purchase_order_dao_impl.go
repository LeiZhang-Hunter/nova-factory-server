package purchasedaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PurchaseOrderDaoImpl 提供 ERP 采购订单数据访问能力。
type PurchaseOrderDaoImpl struct {
	*erpcrud.CRUDDao[purchasemodels.PurchaseOrder, purchasemodels.PurchaseOrderUpsert, purchasemodels.PurchaseOrderQuery]
}

// NewPurchaseOrderDao 创建 ERP 采购订单 DAO。
func NewPurchaseOrderDao(db *gorm.DB) purchasedao.IPurchaseOrderDao {
	return &PurchaseOrderDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[purchasemodels.PurchaseOrder, purchasemodels.PurchaseOrderUpsert, purchasemodels.PurchaseOrderQuery](db, erpcrud.EntityConfig{
			Table:        "erp_purchase_order",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "采购订单号"}},
		}),
	}
}

// List 查询 ERP 采购订单列表。
func (d *PurchaseOrderDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*purchasemodels.PurchaseOrderListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseOrderListData{Rows: result.Rows, Total: result.Total}, nil
}
