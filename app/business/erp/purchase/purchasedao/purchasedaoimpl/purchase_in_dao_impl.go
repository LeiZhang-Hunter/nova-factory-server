package purchasedaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PurchaseInDaoImpl 提供 ERP 采购入库数据访问能力。
type PurchaseInDaoImpl struct {
	*erpcrud.CRUDDao[purchasemodels.PurchaseIn, purchasemodels.PurchaseInUpsert, purchasemodels.PurchaseInQuery]
}

// NewPurchaseInDao 创建 ERP 采购入库 DAO。
func NewPurchaseInDao(db *gorm.DB) purchasedao.IPurchaseInDao {
	return &PurchaseInDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[purchasemodels.PurchaseIn, purchasemodels.PurchaseInUpsert, purchasemodels.PurchaseInQuery](db, erpcrud.EntityConfig{
			Table:        "erp_purchase_in",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "采购入库单号"}},
		}),
	}
}

// List 查询 ERP 采购入库列表。
func (d *PurchaseInDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseInQuery) (*purchasemodels.PurchaseInListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseInListData{Rows: result.Rows, Total: result.Total}, nil
}
