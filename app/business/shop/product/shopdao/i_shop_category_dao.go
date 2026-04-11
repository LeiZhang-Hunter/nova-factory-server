package shopdao

import (
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

type IShopCategoryDao interface {
	Create(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error)
	Update(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopmodels.Category, error)
	List(c *gin.Context, req *shopmodels.CategoryQuery) (*shopmodels.CategoryListData, error)
}
