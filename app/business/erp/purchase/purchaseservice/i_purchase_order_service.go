package purchaseservice

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
)

// IPurchaseOrderService ERP 采购订单服务接口
type IPurchaseOrderService interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error)
	UpdateStatus(c *gin.Context, req *purchasemodels.PurchaseOrderStatusReq) error
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrder, error)
	List(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*purchasemodels.PurchaseOrderListData, error)
}
