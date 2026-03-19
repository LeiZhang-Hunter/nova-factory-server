package shopService

import (
	"nova-factory-server/app/business/shop/shopModels"

	"github.com/gin-gonic/gin"
)

type IShopCategoryService interface {
	Create(c *gin.Context, req *shopModels.CategoryUpsert) (*shopModels.Category, error)
	Update(c *gin.Context, req *shopModels.CategoryUpsert) (*shopModels.Category, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopModels.Category, error)
	List(c *gin.Context, req *shopModels.CategoryQuery) (*shopModels.CategoryListData, error)
}
