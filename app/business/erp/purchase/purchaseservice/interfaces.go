package purchaseservice

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
)

// IPurchaseInService ERP 采购入库服务接口
type IPurchaseInService interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseInUpsert) (*purchasemodels.PurchaseIn, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseInUpsert) (*purchasemodels.PurchaseIn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseIn, error)
	List(c *gin.Context, req *purchasemodels.PurchaseInQuery) (*purchasemodels.PurchaseInListData, error)
}

// IPurchaseInItemService ERP 采购入库项服务接口
type IPurchaseInItemService interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseInItemUpsert) (*purchasemodels.PurchaseInItem, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseInItemUpsert) (*purchasemodels.PurchaseInItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseInItem, error)
	List(c *gin.Context, req *purchasemodels.PurchaseInItemQuery) (*purchasemodels.PurchaseInItemListData, error)
}

// IPurchaseOrderService ERP 采购订单服务接口
type IPurchaseOrderService interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error)
	UpdateStatus(c *gin.Context, req *purchasemodels.PurchaseOrderStatusReq) error
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrder, error)
	List(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*purchasemodels.PurchaseOrderListData, error)
}

// IPurchaseOrderItemService ERP 采购订单项服务接口
type IPurchaseOrderItemService interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert) (*purchasemodels.PurchaseOrderItem, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert) (*purchasemodels.PurchaseOrderItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrderItem, error)
	List(c *gin.Context, req *purchasemodels.PurchaseOrderItemQuery) (*purchasemodels.PurchaseOrderItemListData, error)
}

// IPurchaseReturnService ERP 采购退货服务接口
type IPurchaseReturnService interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseReturn, error)
	List(c *gin.Context, req *purchasemodels.PurchaseReturnQuery) (*purchasemodels.PurchaseReturnListData, error)
}

// IPurchaseReturnItemService ERP 采购退货项服务接口
type IPurchaseReturnItemService interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseReturnItemUpsert) (*purchasemodels.PurchaseReturnItem, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseReturnItemUpsert) (*purchasemodels.PurchaseReturnItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseReturnItem, error)
	List(c *gin.Context, req *purchasemodels.PurchaseReturnItemQuery) (*purchasemodels.PurchaseReturnItemListData, error)
}
