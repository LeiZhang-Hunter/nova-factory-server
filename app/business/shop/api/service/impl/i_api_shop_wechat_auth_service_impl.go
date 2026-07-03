package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	shopusermodels "nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/utils/bCryptPasswordEncoder"
	"strings"
	"time"

	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares/session"

	"github.com/gin-gonic/gin"
)

// IApiShopWechatAuthServiceImpl 提供商城微信授权登录能力。
type IApiShopWechatAuthServiceImpl struct {
	cache     cache.Cache
	configDao dao.IApiShopSysConfigDao
	userDao   dao.IApiShopWechatUserDao
}

// NewIApiShopWechatAuthServiceImpl NewShopWechatAuthService 创建商城微信授权登录服务。
func NewIApiShopWechatAuthServiceImpl(cache cache.Cache, configDao dao.IApiShopSysConfigDao, userDao dao.IApiShopWechatUserDao) service.IApiShopWechatAuthService {
	return &IApiShopWechatAuthServiceImpl{
		cache:     cache,
		configDao: configDao,
		userDao:   userDao,
	}
}

// WechatLogin 微信小程序授权登录
func (s *IApiShopWechatAuthServiceImpl) WechatLogin(c *gin.Context, req *models.WechatLoginReq) (*models.WechatLoginResp, error) {
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
	//appID := "wx64538abe63a51bb3"
	//appSecret := "ed97226ab4c0ef092af0ed6597cbce10"
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

	return s.createLoginSession(c, user)
}

// AccountLogin 使用 shop_user 账号密码登录。
func (s *IApiShopWechatAuthServiceImpl) AccountLogin(c *gin.Context, req *models.AccountLoginReq) (*models.WechatLoginResp, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	account := strings.TrimSpace(req.Account)
	password := strings.TrimSpace(req.Password)
	if account == "" {
		return nil, errors.New("登录账号不能为空")
	}
	if password == "" {
		return nil, errors.New("登录密码不能为空")
	}
	if strings.TrimSpace(req.Code) == "" {
		return nil, errors.New("微信登录凭证不能为空")
	}

	user, err := s.userDao.GetByAccount(c, account)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在或密码错误")
	}
	if !isShopUserEnabled(user.Status) {
		return nil, errors.New("账号已停用")
	}
	if !bCryptPasswordEncoder.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("用户不存在或密码错误")
	}

	openid, err := s.getOpenidByCode(c, strings.TrimSpace(req.Code))
	if err != nil {
		return nil, err
	}
	if err := s.bindAccountOpenid(c, user, openid); err != nil {
		return nil, err
	}

	return s.createLoginSession(c, user)
}

// RefreshToken 刷新 Session Token
func (s *IApiShopWechatAuthServiceImpl) RefreshToken(c *gin.Context, req *models.RefreshTokenReq) (*models.WechatLoginResp, error) {
	if req == nil || req.Token == "" {
		return nil, errors.New("参数不能为空")
	}

	// 通过 session id 获取用户信息
	manager := session.NewShopManager(s.cache)
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
func (s *IApiShopWechatAuthServiceImpl) getWechatConfig(c *gin.Context, key string) (string, error) {
	config, err := s.configDao.GetByConfigKey(c, key)
	if err != nil || config == nil {
		return "", errors.New("配置不存在")
	}
	return config.ConfigValue, nil
}

func (s *IApiShopWechatAuthServiceImpl) getOpenidByCode(c *gin.Context, code string) (string, error) {
	appID, err := s.getWechatConfig(c, "wechat_mini_program_app_id")
	if err != nil {
		return "", errors.New("微信配置缺失")
	}
	appSecret, err := s.getWechatConfig(c, "wechat_mini_program_app_secret")
	if err != nil {
		return "", errors.New("微信配置缺失")
	}
	return s.getWechatOpenid(c, appID, appSecret, code)
}

// getWechatOpenid 调用微信接口获取openid
func (s *IApiShopWechatAuthServiceImpl) getWechatOpenid(c *gin.Context, appID, appSecret, code string) (string, error) {
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
func (s *IApiShopWechatAuthServiceImpl) createWechatUser(c *gin.Context, openid, nickname, avatar string) (*shopusermodels.User, error) {
	// 生成随机用户名
	nickname = fmt.Sprintf("用户-%s", openid[len(openid)-4:])
	userType := int32(1) // 默认代理商类型
	status := true

	userCreate := &models.WechatUserCreate{
		//Username: username,
		Nickname: nickname,
		Avatar:   avatar,
		UserType: userType,
		Status:   &status,
		Openid:   openid,
	}

	return s.userDao.CreateWechatUser(c, userCreate)
}

// getUserDisplayName 获取用户显示名称
func (s *IApiShopWechatAuthServiceImpl) getUserDisplayName(user *shopusermodels.User) string {
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

func (s *IApiShopWechatAuthServiceImpl) bindAccountOpenid(c *gin.Context, user *shopusermodels.User, openid string) error {
	if user == nil {
		return errors.New("用户不存在")
	}
	if user.WechatOpenid != "" {
		if user.WechatOpenid == openid {
			return nil
		}
		return errors.New("当前账号已绑定其他微信")
	}
	boundUser, err := s.userDao.GetByOpenid(c, openid)
	if err != nil {
		return errors.New("查询微信绑定状态失败")
	}
	if boundUser != nil && boundUser.ID != user.ID {
		return errors.New("该微信已绑定其他账号")
	}
	if err := s.userDao.BindWechatOpenid(c, user.ID, openid); err != nil {
		return errors.New("绑定微信失败")
	}
	user.WechatOpenid = openid
	return nil
}

func (s *IApiShopWechatAuthServiceImpl) createLoginSession(c *gin.Context, user *shopusermodels.User) (*models.WechatLoginResp, error) {
	manager := session.NewShopManager(s.cache)
	currentSession, err := manager.InitSessionWithData(c, user.ID, &session.SessionData{
		SessionType: sessionStatus.SessionTypeShopUser,
		UserId:      user.ID,
		DeptId:      user.DeptID,
		Avatar:      user.Avatar,
		UserName:    s.getUserDisplayName(user),
		IpAddr:      c.ClientIP(),
		LoginTime:   time.Now().Unix(),
	})
	if err != nil {
		return nil, errors.New("创建会话失败")
	}
	return &models.WechatLoginResp{
		Token:  currentSession.Id(),
		UserId: user.ID,
	}, nil
}

func isShopUserEnabled(status *bool) bool {
	return status != nil && *status
}
