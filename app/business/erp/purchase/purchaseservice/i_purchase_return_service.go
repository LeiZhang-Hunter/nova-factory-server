package purchaseservice

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
)

// IPurchaseReturnService ERP 采购退货服务接口
type IPurchaseReturnService interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseReturn, error)
	List(c *gin.Context, req *purchasemodels.PurchaseReturnQuery) (*purchasemodels.PurchaseReturnListData, error)
}
