package dao

import (
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// IShopCartDao 商城用户购物车数据访问接口
type IShopCartDao interface {
	Set(c *gin.Context, req *models.CartSetReq) (*models.Cart, error)
	GetByID(c *gin.Context, id int64) (*models.Cart, error)
	GetByUserIDAndSkuID(c *gin.Context, userID int64, skuID string) (*models.Cart, error)
	List(c *gin.Context, req *models.CartQuery) (*models.CartListData, error)
	Remove(c *gin.Context, ids []int64) error
}
