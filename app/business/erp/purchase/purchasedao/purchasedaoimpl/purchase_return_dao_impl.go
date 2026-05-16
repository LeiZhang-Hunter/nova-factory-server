package purchasedaoimpl

import (
	"nova-factory-server/app/business/erp/erpcrud"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PurchaseReturnDaoImpl 提供 ERP 采购退货数据访问能力。
type PurchaseReturnDaoImpl struct {
	*erpcrud.CRUDDao[purchasemodels.PurchaseReturn, purchasemodels.PurchaseReturnUpsert, purchasemodels.PurchaseReturnQuery]
}

// NewPurchaseReturnDao 创建 ERP 采购退货 DAO。
func NewPurchaseReturnDao(db *gorm.DB) purchasedao.IPurchaseReturnDao {
	return &PurchaseReturnDaoImpl{
		CRUDDao: erpcrud.NewCRUDDao[purchasemodels.PurchaseReturn, purchasemodels.PurchaseReturnUpsert, purchasemodels.PurchaseReturnQuery](db, erpcrud.EntityConfig{
			Table:        "erp_purchase_return",
			OrderBy:      "id DESC",
			UniqueFields: []erpcrud.UniqueField{{Field: "No", Column: "no", Label: "采购退货单号"}},
		}),
	}
}

// List 查询 ERP 采购退货列表。
func (d *PurchaseReturnDaoImpl) List(c *gin.Context, req *purchasemodels.PurchaseReturnQuery) (*purchasemodels.PurchaseReturnListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseReturnListData{Rows: result.Rows, Total: result.Total}, nil
}
