package impl

import (
	"errors"

	"nova-factory-server/app/business/shop/api/dao"
	apiModels "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// ShopAuthServiceImpl 提供商城小程序鉴权服务。
type ShopAuthServiceImpl struct {
	cache   cache.Cache
	userDao dao.IShopWechatUserDao
}

// NewShopAuthService 创建商城小程序鉴权服务。
func NewShopAuthService(cache cache.Cache, userDao dao.IShopWechatUserDao) service.IAppShopAuthService {
	return &ShopAuthServiceImpl{
		cache:   cache,
		userDao: userDao,
	}
}

// GetInfo 获取当前商城登录用户信息（小程序）。
func (s *ShopAuthServiceImpl) GetInfo(c *gin.Context) (*apiModels.ShopGetInfoResp, error) {
	userID := baizeContext.GetUserId(c)
	if userID == 0 {
		return nil, errors.New("用户未登录")
	}

	user, err := s.userDao.GetByUserID(c, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	return &apiModels.ShopGetInfoResp{
		User: apiModels.UserToAuthUserInfo(user),
	}, nil
}

// Logout 退出当前商城登录会话（小程序）。
func (s *ShopAuthServiceImpl) Logout(c *gin.Context) error {
	session.NewManger(s.cache).RemoveSession(c)
	return nil
}
