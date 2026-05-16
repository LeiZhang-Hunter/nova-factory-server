package masterservice

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// ISupplierService ERP 供应商服务接口
type ISupplierService interface {
	Create(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error)
	Update(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Supplier, error)
	List(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error)
}
