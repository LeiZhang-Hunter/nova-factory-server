package service

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IAppShopWechatAuthService  商城微信授权登录服务接口
type IAppShopWechatAuthService interface {
	WechatLogin(c *gin.Context, req *models.WechatLoginReq) (*models.WechatLoginResp, error)
	RefreshToken(c *gin.Context, req *models.RefreshTokenReq) (*models.WechatLoginResp, error)
}
