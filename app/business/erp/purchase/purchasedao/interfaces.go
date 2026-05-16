package purchasedao

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpcrud"
)

// IPurchaseInDao ERP 采购入库数据访问接口
type IPurchaseInDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseInUpsert) (*purchasemodels.PurchaseIn, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseInUpsert) (*purchasemodels.PurchaseIn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseIn, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseIn, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseInQuery) (*erpcrud.PageResult[purchasemodels.PurchaseIn], error)
	List(c *gin.Context, req *purchasemodels.PurchaseInQuery) (*purchasemodels.PurchaseInListData, error)
}

// IPurchaseInItemDao ERP 采购入库项数据访问接口
type IPurchaseInItemDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseInItemUpsert) (*purchasemodels.PurchaseInItem, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseInItemUpsert) (*purchasemodels.PurchaseInItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseInItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseInItem, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseInItemQuery) (*erpcrud.PageResult[purchasemodels.PurchaseInItem], error)
	List(c *gin.Context, req *purchasemodels.PurchaseInItemQuery) (*purchasemodels.PurchaseInItemListData, error)
}

// IPurchaseOrderDao ERP 采购订单数据访问接口
type IPurchaseOrderDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseOrderUpsert) (*purchasemodels.PurchaseOrder, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrder, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseOrder, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*erpcrud.PageResult[purchasemodels.PurchaseOrder], error)
	List(c *gin.Context, req *purchasemodels.PurchaseOrderQuery) (*purchasemodels.PurchaseOrderListData, error)
}

// IPurchaseOrderItemDao ERP 采购订单项数据访问接口
type IPurchaseOrderItemDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert) (*purchasemodels.PurchaseOrderItem, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert) (*purchasemodels.PurchaseOrderItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrderItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseOrderItem, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseOrderItemQuery) (*erpcrud.PageResult[purchasemodels.PurchaseOrderItem], error)
	List(c *gin.Context, req *purchasemodels.PurchaseOrderItemQuery) (*purchasemodels.PurchaseOrderItemListData, error)
}

// IPurchaseReturnDao ERP 采购退货数据访问接口
type IPurchaseReturnDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseReturn, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseReturn, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseReturnQuery) (*erpcrud.PageResult[purchasemodels.PurchaseReturn], error)
	List(c *gin.Context, req *purchasemodels.PurchaseReturnQuery) (*purchasemodels.PurchaseReturnListData, error)
}

// IPurchaseReturnItemDao ERP 采购退货项数据访问接口
type IPurchaseReturnItemDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseReturnItemUpsert) (*purchasemodels.PurchaseReturnItem, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseReturnItemUpsert) (*purchasemodels.PurchaseReturnItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseReturnItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseReturnItem, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseReturnItemQuery) (*erpcrud.PageResult[purchasemodels.PurchaseReturnItem], error)
	List(c *gin.Context, req *purchasemodels.PurchaseReturnItemQuery) (*purchasemodels.PurchaseReturnItemListData, error)
}
