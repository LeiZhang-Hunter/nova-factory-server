package purchaseservice

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
)

// IPurchaseInItemService ERP 采购入库项服务接口
type IPurchaseInItemService interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseInItemUpsert) (*purchasemodels.PurchaseInItem, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseInItemUpsert) (*purchasemodels.PurchaseInItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseInItem, error)
	List(c *gin.Context, req *purchasemodels.PurchaseInItemQuery) (*purchasemodels.PurchaseInItemListData, error)
}
