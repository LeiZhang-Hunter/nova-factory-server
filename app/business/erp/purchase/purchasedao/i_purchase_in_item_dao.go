package purchasedao

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IPurchaseInItemDao ERP 采购入库项数据访问接口
type IPurchaseInItemDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseInItemUpsert) (*purchasemodels.PurchaseInItem, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseInItemUpsert) (*purchasemodels.PurchaseInItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseInItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseInItem, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseInItemQuery) (*erpbiz.PageResult[purchasemodels.PurchaseInItem], error)
	List(c *gin.Context, req *purchasemodels.PurchaseInItemQuery) (*purchasemodels.PurchaseInItemListData, error)
}
