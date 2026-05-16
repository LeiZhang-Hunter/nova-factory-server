package purchasedao

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IPurchaseReturnItemDao ERP 采购退货项数据访问接口
type IPurchaseReturnItemDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseReturnItemUpsert) (*purchasemodels.PurchaseReturnItem, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseReturnItemUpsert) (*purchasemodels.PurchaseReturnItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseReturnItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseReturnItem, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseReturnItemQuery) (*erpbiz.PageResult[purchasemodels.PurchaseReturnItem], error)
	List(c *gin.Context, req *purchasemodels.PurchaseReturnItemQuery) (*purchasemodels.PurchaseReturnItemListData, error)
}
