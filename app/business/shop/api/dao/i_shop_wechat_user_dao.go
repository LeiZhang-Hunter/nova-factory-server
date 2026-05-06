package dao

import (
	wechatModels "nova-factory-server/app/business/shop/api/models"
	userModels "nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// IShopWechatUserDao 商城微信用户数据访问接口
type IShopWechatUserDao interface {
	GetByOpenid(c *gin.Context, openid string) (*userModels.User, error)
	CreateWechatUser(c *gin.Context, req *wechatModels.WechatUserCreate) (*userModels.User, error)
}
