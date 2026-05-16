package purchasedaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PurchaseInItemDaoImpl 提供 ERP 采购入库项数据访问能力。
type PurchaseInItemDaoImpl struct {
	*erpcrud.CRUDDao[purchasemodels.PurchaseInItem, purchasemodels.PurchaseInItemUpsert, purchasemodels.PurchaseInItemQuery]
}

// NewPurchaseInItemDao 创建 ERP 采购入库项 DAO。
func NewPurchaseInItemDao(db *gorm.DB) purchasedao.IPurchaseInItemDao {
	return &PurchaseInItemDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[purchasemodels.PurchaseInItem, purchasemodels.PurchaseInItemUpsert, purchasemodels.PurchaseInItemQuery](db, erpcrud.EntityConfig{
			Table:        "erp_purchase_in_items",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{},
		}),
	}
}

// List 查询 ERP 采购入库项列表。
func (d *PurchaseInItemDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseInItemQuery) (*purchasemodels.PurchaseInItemListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseInItemListData{Rows: result.Rows, Total: result.Total}, nil
}
