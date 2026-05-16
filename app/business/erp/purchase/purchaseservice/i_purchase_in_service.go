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
