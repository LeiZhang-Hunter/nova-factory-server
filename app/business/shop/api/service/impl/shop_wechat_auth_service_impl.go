package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	configDao "nova-factory-server/app/business/shop/config/dao"
	userModels "nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"

	"github.com/gin-gonic/gin"
)

// ShopWechatAuthServiceImpl 提供商城微信授权登录能力。
type ShopWechatAuthServiceImpl struct {
	cache     cache.Cache
	configDao configDao.IShopSysConfigDao
	userDao   dao.IShopWechatUserDao
}

// NewShopWechatAuthService 创建商城微信授权登录服务。
func NewShopWechatAuthService(cache cache.Cache, configDao configDao.IShopSysConfigDao, userDao dao.IShopWechatUserDao) service.IShopWechatAuthService {
	return &ShopWechatAuthServiceImpl{
		cache:     cache,
		configDao: configDao,
		userDao:   userDao,
	}
}

// WechatLogin 微信小程序授权登录
func (s *ShopWechatAuthServiceImpl) WechatLogin(c *gin.Context, req *models.WechatLoginReq) (*models.WechatLoginResp, error) {
	if req == nil || req.Code == "" {
		return nil, errors.New("参数不能为空")
	}

	// 获取微信配置
	appID, err := s.getWechatConfig(c, "wechat_mini_program_app_id")
	if err != nil {
		return nil, errors.New("微信配置缺失")
	}
	appSecret, err := s.getWechatConfig(c, "wechat_mini_program_app_secret")
	if err != nil {
		return nil, errors.New("微信配置缺失")
	}

	// 调用微信接口获取openid
	openid, err := s.getWechatOpenid(c, appID, appSecret, req.Code)
	if err != nil {
		return nil, err
	}

	// 查找是否已存在绑定该openid的用户
	user, err := s.userDao.GetByOpenid(c, openid)
	if err != nil {
		return nil, errors.New("查询用户失败")
	}

	// 未绑定则创建新用户
	if user == nil {
		user, err = s.createWechatUser(c, openid, req.Nickname, req.Avatar)
		if err != nil {
			return nil, err
		}
	}

	// 创建 Session
	manager := session.NewManger(s.cache)
	currentSession, err := manager.InitSession(c, user.ID)
	if err != nil {
		return nil, errors.New("创建会话失败")
	}
	currentSession.Set(c, sessionStatus.SessionType, sessionStatus.SessionTypeShopUser)
	currentSession.Set(c, sessionStatus.UserId, user.ID)
	currentSession.Set(c, sessionStatus.UserName, s.getUserDisplayName(user))
	currentSession.Set(c, sessionStatus.Avatar, user.Avatar)

	return &models.WechatLoginResp{
		Token:  currentSession.Id(),
		UserId: user.ID,
	}, nil
}

// RefreshToken 刷新 Session Token
func (s *ShopWechatAuthServiceImpl) RefreshToken(c *gin.Context, req *models.RefreshTokenReq) (*models.WechatLoginResp, error) {
	if req == nil || req.Token == "" {
		return nil, errors.New("参数不能为空")
	}

	// 通过 session id 获取用户信息
	manager := session.NewManger(s.cache)
	sess, err := manager.Get(c, req.Token)
	if err != nil {
		return nil, errors.New("无效的会话")
	}

	userIdStr := sess.Get(c, sessionStatus.UserId)
	var userId int64
	fmt.Sscanf(userIdStr, "%d", &userId)

	return &models.WechatLoginResp{
		Token:  sess.Id(),
		UserId: userId,
	}, nil
}

// getWechatConfig 获取微信配置
func (s *ShopWechatAuthServiceImpl) getWechatConfig(c *gin.Context, key string) (string, error) {
	config, err := s.configDao.GetByConfigKey(c, key)
	if err != nil || config == nil {
		return "", errors.New("配置不存在")
	}
	return config.ConfigValue, nil
}

// getWechatOpenid 调用微信接口获取openid
func (s *ShopWechatAuthServiceImpl) getWechatOpenid(c *gin.Context, appID, appSecret, code string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appID, appSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("微信接口调用失败")
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", errors.New("解析微信响应失败")
	}

	if errMsg, ok := result["errmsg"]; ok {
		return "", fmt.Errorf("微信错误: %v", errMsg)
	}

	openid, ok := result["openid"].(string)
	if !ok || openid == "" {
		return "", errors.New("获取openid失败")
	}

	return openid, nil
}

// createWechatUser 创建微信用户
func (s *ShopWechatAuthServiceImpl) createWechatUser(c *gin.Context, openid, nickname, avatar string) (*userModels.User, error) {
	// 生成随机用户名
	username := fmt.Sprintf("wx_%s", openid[:16])
	userType := int32(1) // 默认代理商类型
	status := true

	userCreate := &models.WechatUserCreate{
		Username: username,
		Nickname: nickname,
		Avatar:   avatar,
		UserType: userType,
		Status:   &status,
		Openid:   openid,
	}

	return s.userDao.CreateWechatUser(c, userCreate)
}

// getUserDisplayName 获取用户显示名称
func (s *ShopWechatAuthServiceImpl) getUserDisplayName(user *userModels.User) string {
	if user == nil {
		return ""
	}
	if user.Nickname != "" {
		return user.Nickname
	}
	if user.ContactName != "" {
		return user.ContactName
	}
	return user.Username
}
