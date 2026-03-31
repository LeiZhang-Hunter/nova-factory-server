package shopservice

import (
	"nova-factory-server/app/business/shop/shopmodels"

	"github.com/gin-gonic/gin"
)

type IShopUserService interface {
	Create(c *gin.Context, req *shopmodels.UserUpsert) (*shopmodels.User, error)
	Update(c *gin.Context, req *shopmodels.UserUpsert) (*shopmodels.User, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*shopmodels.User, error)
	List(c *gin.Context, req *shopmodels.UserQuery) (*shopmodels.UserListData, error)
}
