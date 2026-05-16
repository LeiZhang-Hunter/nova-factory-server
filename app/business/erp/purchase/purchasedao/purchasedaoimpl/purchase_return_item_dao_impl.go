package purchasedaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PurchaseReturnItemDaoImpl 提供 ERP 采购退货项数据访问能力。
type PurchaseReturnItemDaoImpl struct {
	*erpcrud.CRUDDao[purchasemodels.PurchaseReturnItem, purchasemodels.PurchaseReturnItemUpsert, purchasemodels.PurchaseReturnItemQuery]
}

// NewPurchaseReturnItemDao 创建 ERP 采购退货项 DAO。
func NewPurchaseReturnItemDao(db *gorm.DB) purchasedao.IPurchaseReturnItemDao {
	return &PurchaseReturnItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[purchasemodels.PurchaseReturnItem, purchasemodels.PurchaseReturnItemUpsert, purchasemodels.PurchaseReturnItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_purchase_return_items",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 采购退货项列表。
func (d *PurchaseReturnItemDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseReturnItemQuery) (*purchasemodels.PurchaseReturnItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseReturnItemListData{Rows: result.Rows, Total: result.Total}, nil
}
