package shopDao

import (
	"nova-factory-server/app/business/shop/shopModels"

	"github.com/gin-gonic/gin"
)

type IShopUserDao interface {
	Create(c *gin.Context, req *shopModels.UserUpsert) (*shopModels.User, error)
	Update(c *gin.Context, req *shopModels.UserUpsert) (*shopModels.User, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopModels.User, error)
	List(c *gin.Context, req *shopModels.UserQuery) (*shopModels.UserListData, error)
}
