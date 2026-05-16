package masterdao

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IProductDao ERP 产品数据访问接口
type IProductDao interface {
	Create(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error)
	Update(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Product, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.Product, error)
	ListPage(c *gin.Context, req *mastermodels.ProductQuery) (*erpbiz.PageResult[mastermodels.Product], error)
	List(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error)
}
