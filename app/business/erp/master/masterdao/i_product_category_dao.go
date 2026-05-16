package masterdao

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IProductCategoryDao ERP 产品分类数据访问接口
type IProductCategoryDao interface {
	Create(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error)
	Update(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.ProductCategory, error)
	GetByColumn(c *gin.Context, column string, value any) (*mastermodels.ProductCategory, error)
	ListPage(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*erpbiz.PageResult[mastermodels.ProductCategory], error)
	List(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error)
}
