package masterdao

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// ISupplierDao ERP 供应商数据访问接口
type ISupplierDao interface {
	Create(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error)
	Update(c *gin.Context, req *mastermodels.SupplierUpsert) (*mastermodels.Supplier, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Supplier, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Supplier, error)
	ListPage(c *gin.Context, req *mastermodels.SupplierQuery) (*erpbiz.PageResult[mastermodels.Supplier], error)
	List(c *gin.Context, req *mastermodels.SupplierQuery) (*mastermodels.SupplierListData, error)
}
