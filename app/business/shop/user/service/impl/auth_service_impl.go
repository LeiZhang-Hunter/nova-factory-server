package impl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"
	"nova-factory-server/app/utils/bCryptPasswordEncoder"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// ShopAuthServiceImpl 提供商城用户登录与鉴权能力。
type ShopAuthServiceImpl struct {
	cache   cache.Cache
	userDao dao.IShopUserDao
}

// NewShopAuthService 创建商城用户鉴权服务。
func NewShopAuthService(cache cache.Cache, userDao dao.IShopUserDao) service.IShopAuthService {
	return &ShopAuthServiceImpl{
		cache:   cache,
		userDao: userDao,
	}
}

// Login 使用 shop_user 账号登录并初始化商城会话。
func (s *ShopAuthServiceImpl) Login(c *gin.Context, req *models.ShopLoginReq) (*models.ShopLoginResp, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	req.Account = strings.TrimSpace(req.Account)
	if req.Account == "" {
		return nil, errors.New("登录账号不能为空")
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, errors.New("登录密码不能为空")
	}

	user, err := s.userDao.GetByAccount(c, req.Account)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在/密码错误")
	}
	if !isShopUserEnabled(user.Status) {
		return nil, errors.New("账号已停用")
	}
	if !bCryptPasswordEncoder.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("用户不存在/密码错误")
	}

	manager := session.NewManger(s.cache)
	currentSession, err := manager.InitSession(c, user.ID)
	if err != nil {
		return nil, err
	}
	permissions := shopPermissionsByUser(user)
	currentSession.Set(c, sessionStatus.SessionType, sessionStatus.SessionTypeShopUser)
	currentSession.Set(c, sessionStatus.UserId, user.ID)
	currentSession.Set(c, sessionStatus.UserName, shopUserDisplayName(user))
	currentSession.Set(c, sessionStatus.Avatar, user.Avatar)
	currentSession.Set(c, sessionStatus.DeptId, user.DeptID)
	currentSession.Set(c, sessionStatus.Permission, permissions)

	return &models.ShopLoginResp{
		Token:       currentSession.Id(),
		TokenType:   "Bearer",
		User:        buildShopAuthUserInfo(user),
		Permissions: permissions,
	}, nil
}

// GetInfo 获取当前商城登录用户信息。
func (s *ShopAuthServiceImpl) GetInfo(c *gin.Context) (*models.ShopGetInfoResp, error) {
	userID := baizeContext.GetUserId(c)
	user, err := s.userDao.GetByID(c, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("商城用户不存在")
	}
	permissions := shopPermissionsByUser(user)
	sessionValue := baizeContext.GetSession(c)
	sessionValue.Set(c, sessionStatus.UserName, shopUserDisplayName(user))
	sessionValue.Set(c, sessionStatus.Avatar, user.Avatar)
	sessionValue.Set(c, sessionStatus.DeptId, user.DeptID)
	sessionValue.Set(c, sessionStatus.Permission, permissions)
	return &models.ShopGetInfoResp{
		User:        buildShopAuthUserInfo(user),
		Permissions: permissions,
	}, nil
}

// Logout 退出当前商城登录会话。
func (s *ShopAuthServiceImpl) Logout(c *gin.Context) error {
	session.NewManger(s.cache).RemoveSession(c)
	return nil
}

func buildShopAuthUserInfo(user *models.User) *models.ShopAuthUserInfo {
	if user == nil {
		return nil
	}
	return &models.ShopAuthUserInfo{
		ID:           user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Mobile:       user.Mobile,
		Avatar:       user.Avatar,
		UserType:     user.UserType,
		CompanyName:  user.CompanyName,
		ContactName:  user.ContactName,
		ContactPhone: user.ContactPhone,
		DeptID:       user.DeptID,
	}
}

func shopPermissionsByUser(user *models.User) []string {
	permissions := []string{"shop:app:user"}
	if user == nil {
		return permissions
	}
	switch user.UserType {
	case 1:
		permissions = append(permissions, "shop:app:agent")
	case 2:
		permissions = append(permissions, "shop:app:distributor")
	case 3:
		permissions = append(permissions, "shop:app:factory")
	}
	return permissions
}

func shopUserDisplayName(user *models.User) string {
	if user == nil {
		return ""
	}
	if value := strings.TrimSpace(user.Nickname); value != "" {
		return value
	}
	if value := strings.TrimSpace(user.ContactName); value != "" {
		return value
	}
	return strings.TrimSpace(user.Username)
}

func isShopUserEnabled(status *bool) bool {
	return status != nil && *status
}
