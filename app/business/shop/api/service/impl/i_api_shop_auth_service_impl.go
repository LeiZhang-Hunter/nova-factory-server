package impl

import (
	"errors"

	"nova-factory-server/app/business/shop/api/dao"
	apiModels "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
	"nova-factory-server/app/utils/bCryptPasswordEncoder"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// IApiShopAuthServiceImpl 提供商城小程序鉴权服务。
type IApiShopAuthServiceImpl struct {
	cache   cache.Cache
	userDao dao.IApiShopWechatUserDao
}

// NewIApiShopAuthServiceImpl   创建商城小程序鉴权服务。
func NewIApiShopAuthServiceImpl(cache cache.Cache, userDao dao.IApiShopWechatUserDao) service.IApiShopAuthService {
	return &IApiShopAuthServiceImpl{
		cache:   cache,
		userDao: userDao,
	}
}

// GetInfo 获取当前商城登录用户信息（小程序）。
func (s *IApiShopAuthServiceImpl) GetInfo(c *gin.Context) (*apiModels.ShopGetInfoResp, error) {
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
func (s *IApiShopAuthServiceImpl) Logout(c *gin.Context) error {
	session.NewShopManager(s.cache).RemoveSession(c)
	return nil
}

// ChangePassword 修改商城用户密码。
func (s *IApiShopAuthServiceImpl) ChangePassword(c *gin.Context, req *apiModels.ChangePasswordReq) error {
	userID := baizeContext.GetUserId(c)
	if userID == 0 {
		return errors.New("用户未登录")
	}

	user, err := s.userDao.GetByUserID(c, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	if req.OldPassword == req.NewPassword {
		return errors.New("新密码不能与旧密码相同")
	}

	if len(req.NewPassword) < 6 {
		return errors.New("新密码长度不能少于6位")
	}

	if !bCryptPasswordEncoder.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.New("旧密码错误")
	}

	hashedPassword := bCryptPasswordEncoder.HashPassword(req.NewPassword)
	return s.userDao.UpdatePassword(c, userID, hashedPassword)
}

// UpdateProfile 更新商城用户个人资料。
func (s *IApiShopAuthServiceImpl) UpdateProfile(c *gin.Context, req *apiModels.UpdateProfileReq) error {
	userID := baizeContext.GetUserId(c)
	if userID == 0 {
		return errors.New("用户未登录")
	}

	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Mobile != "" {
		updates["mobile"] = req.Mobile
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	if len(updates) == 0 {
		return errors.New("没有需要更新的字段")
	}

	return s.userDao.UpdateProfile(c, userID, updates)
}
