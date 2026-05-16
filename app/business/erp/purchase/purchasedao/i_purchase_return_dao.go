package purchasedao

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
)

// IPurchaseReturnDao ERP 采购退货数据访问接口
type IPurchaseReturnDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseReturn, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseReturn, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseReturnQuery) (*purchasemodels.PurchaseReturnListData, error)
	List(c *gin.Context, req *purchasemodels.PurchaseReturnQuery) (*purchasemodels.PurchaseReturnListData, error)
}
