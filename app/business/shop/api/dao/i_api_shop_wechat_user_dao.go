package dao

import (
	"nova-factory-server/app/business/shop/api/models"
	shopusermodels "nova-factory-server/app/business/shop/user/models"

	"github.com/gin-gonic/gin"
)

// IApiShopWechatUserDao  商城微信用户数据访问接口
type IApiShopWechatUserDao interface {
	GetByOpenid(c *gin.Context, openid string) (*shopusermodels.User, error)
	CreateWechatUser(c *gin.Context, req *models.WechatUserCreate) (*shopusermodels.User, error)
	GetByID(c *gin.Context, id int64) (*shopusermodels.User, error)
	GetByUserID(c *gin.Context, userId int64) (*shopusermodels.User, error)
}
