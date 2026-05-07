package service

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IAppShopAuthService  商城小程序鉴权服务接口
type IAppShopAuthService interface {
	GetInfo(c *gin.Context) (*models.ShopGetInfoResp, error)
	Logout(c *gin.Context) error
}
