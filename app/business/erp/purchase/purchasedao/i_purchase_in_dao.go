package purchasedao

import (
	"nova-factory-server/app/business/erp/purchase/purchasemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IPurchaseInDao ERP 采购入库数据访问接口
type IPurchaseInDao interface {
	Create(c *gin.Context, req *purchasemodels.PurchaseInUpsert) (*purchasemodels.PurchaseIn, error)
	Update(c *gin.Context, req *purchasemodels.PurchaseInUpsert) (*purchasemodels.PurchaseIn, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseIn, error)
	GetByColumn(c *gin.Context, column string, value any) (*purchasemodels.PurchaseIn, error)
	ListPage(c *gin.Context, req *purchasemodels.PurchaseInQuery) (*erpbiz.PageResult[purchasemodels.PurchaseIn], error)
	List(c *gin.Context, req *purchasemodels.PurchaseInQuery) (*purchasemodels.PurchaseInListData, error)
}
