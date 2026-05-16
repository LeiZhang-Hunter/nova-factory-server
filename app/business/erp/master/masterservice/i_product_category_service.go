package masterservice

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// IProductCategoryService ERP 产品分类服务接口
type IProductCategoryService interface {
	Create(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error)
	Update(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.ProductCategory, error)
	List(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error)
}
