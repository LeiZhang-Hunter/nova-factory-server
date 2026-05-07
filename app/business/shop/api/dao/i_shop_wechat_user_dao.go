package dao

import (
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

// IShopWechatUserDao 商城微信用户数据访问接口
type IShopWechatUserDao interface {
	GetByOpenid(c *gin.Context, openid string) (*models.User, error)
	CreateWechatUser(c *gin.Context, req *models.WechatUserCreate) (*models.User, error)
	GetByID(c *gin.Context, id int64) (*models.User, error)
	GetByUserID(c *gin.Context, id int64) (*models.User, error)
}
