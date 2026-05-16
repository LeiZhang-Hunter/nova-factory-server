package purchasedao

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
)

// IPurchaseOrderDao ERP 采购订单数据访问接口
type IPurchaseOrderDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrder, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseOrder, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*purchasemodels.PurchaseOrderListData, error)
	List(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*purchasemodels.PurchaseOrderListData, error)
}
