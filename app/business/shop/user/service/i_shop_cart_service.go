package service

import (
	api "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// IShopCartService 商城用户购物车服务接口
type IShopCartService interface {
	Set(c *gin.Context, req *models.CartSetReq) (*models.Cart, error)
	GetByID(c *gin.Context, id int64) (*models.Cart, error)
	List(c *gin.Context, req *models.CartQuery) (*models.CartListData, error)
	Remove(c *gin.Context, ids []int64) error
}
