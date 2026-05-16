package purchasedaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PurchaseOrderItemDaoImpl 提供 ERP 采购订单项数据访问能力。
type PurchaseOrderItemDaoImpl struct {
	*erpcrud.CRUDDao[purchasemodels.PurchaseOrderItem, purchasemodels.PurchaseOrderItemUpsert, purchasemodels.PurchaseOrderItemQuery]
}

// NewPurchaseOrderItemDao 创建 ERP 采购订单项 DAO。
func NewPurchaseOrderItemDao(db *gorm.DB) purchasedao.IPurchaseOrderItemDao {
	return &PurchaseOrderItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[purchasemodels.PurchaseOrderItem, purchasemodels.PurchaseOrderItemUpsert, purchasemodels.PurchaseOrderItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_purchase_order_items",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 采购订单项列表。
func (d *PurchaseOrderItemDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseOrderItemQuery) (*purchasemodels.PurchaseOrderItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseOrderItemListData{Rows: result.Rows, Total: result.Total}, nil
}
