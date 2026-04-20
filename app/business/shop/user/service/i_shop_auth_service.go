package service

import (
	"nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// IShopAuthService 商城用户登录与鉴权服务接口
type IShopAuthService interface {
	Login(c *gin.Context, req *models.ShopLoginReq) (*models.ShopLoginResp, error)
	GetInfo(c *gin.Context) (*models.ShopGetInfoResp, error)
	Logout(c *gin.Context) error
}
