package dao

import (
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

type IShopUserDao interface {
	Create(c *gin.Context, req *models.UserUpsert) (*models.User, error)
	Update(c *gin.Context, req *models.UserUpsert) (*models.User, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*models.User, error)
	List(c *gin.Context, req *models.UserQuery) (*models.UserListData, error)
}
