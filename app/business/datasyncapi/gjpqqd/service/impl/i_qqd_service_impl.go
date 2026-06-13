// 管家婆全渠道服务实现。
// 实现 service.GjpQqdService 接口，负责 OAuth 授权码生成、
// Token 签发/刷新/校验、MD5 签名校验及商品/库存操作。
// 所有 Token 与授权码通过缓存存储，支持 TTL 过期自动失效。
package impl

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"nova-factory-server/app/business/datasyncapi/gjpqqd/models"
	"nova-factory-server/app/business/datasyncapi/gjpqqd/service"
	"nova-factory-server/app/business/erp/setting/settingdao"
	"nova-factory-server/app/datasource/cache"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// cacheKeyPrefix 缓存键前缀，所有管家婆全渠道相关的缓存键均以此为前缀
const cacheKeyPrefix = "erp_api:qqd:"

// IQQDServiceImpl 管家婆全渠道服务实现
// 持有集成配置缓存、配置 DAO（用于初始化配置）
// goodsDao/goodsSkuDao 预留用于未来的数据库直接操作
type IQQDServiceImpl struct {
	cfg   models.QQDConfig
	cache cache.Cache
	// goodsDao    预留：商品 DAO
	// goodsSkuDao 预留：商品 SKU DAO
	configDao settingdao.IIntegrationConfigDao
}

func (s *IQQDServiceImpl) ProductList(ctx *gin.Context, request models.ProductListRequest) (map[string]any, error) {
	//TODO implement me
	panic("implement me")
}

func (s *IQQDServiceImpl) AddProducts(ctx *gin.Context, goodsInfos []map[string]any) ([]map[string]any, error) {
	//TODO implement me
	panic("implement me")
}

func (s *IQQDServiceImpl) ProductStockUpdate(ctx *gin.Context, request models.ProductStockUpdateRequest) (map[string]any, error) {
	//TODO implement me
	panic("implement me")
}

// authCodePayload 授权码的缓存负载，记录授权的 appKey、回调地址和过期时间
type authCodePayload struct {
	AppKey      string    `json:"app_key"`
	RedirectURI string    `json:"redirect_uri"`
	State       string    `json:"state"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// tokenPayload access_token 或 refresh_token 的缓存负载
type tokenPayload struct {
	AppKey    string    `json:"app_key"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewIQQDServiceImpl 创建管家婆全渠道服务实例
func NewIQQDServiceImpl(configDao settingdao.IIntegrationConfigDao, cache cache.Cache) service.GjpQqdService {
	return &IQQDServiceImpl{
		configDao: configDao,
		cache:     cache,
	}
}

// initConfig 从集成配置表加载管家婆配置到 s.cfg
// 每次业务操作前调用，确保使用最新配置
func (s *IQQDServiceImpl) initConfig(ctx *gin.Context) error {
	enabled, err := s.configDao.GetEnabled(ctx)
	if err != nil {
		zap.L().Error("get integration config fail", zap.Error(err))
		return err
	}
	if enabled == nil {
		return errors.New("integration config disabled")
	}
	enableService, err := enabled.Service()
	if err != nil {
		return err
	}
	if enableService == nil {
		return errors.New("service disabled")
	}
	var cfg models.QQDConfig
	cfg.ApplyDefaults()
	err = json.Unmarshal([]byte(enabled.Data), &cfg)
	if err != nil {
		zap.L().Error("parse qqd config fail", zap.Error(err))
		return err
	}
	s.cfg = cfg
	return nil
}

// CreateAuthorizationCallback 校验 app 凭据，生成一次性授权码存入缓存
// 返回拼接好 code 和 state 的 redirect_uri，供浏览器 302 跳转
func (s *IQQDServiceImpl) CreateAuthorizationCallback(ctx *gin.Context, req *models.AuthorizeReq) (string, error) {
	if s.initConfig(ctx) != nil {
		return "", models.ErrInvalidCredential
	}
	if !s.validCredential(req.AppKey, req.AppSecret) {
		return "", models.ErrInvalidCredential
	}
	code, err := randomToken()
	if err != nil {
		return "", fmt.Errorf("generate auth code: %w", err)
	}
	codeTTL := parseDurationOrDefault(s.cfg.CodeTTL, 10*time.Minute)
	payload := authCodePayload{
		AppKey:      req.AppKey,
		RedirectURI: req.Redirect,
		State:       req.State,
		ExpiresAt:   time.Now().Add(codeTTL),
	}
	// 将授权码存入缓存，设置 TTL
	if err := s.setJSON(ctx, authCodeKey(code), payload, codeTTL); err != nil {
		return "", fmt.Errorf("save auth code: %w", err)
	}
	parsed, err := url.Parse(req.Redirect)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	query.Set("code", code)
	query.Set("state", req.State)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

// IssueToken 签发或刷新 access_token
// - authorization_code 模式：消耗一次性授权码，返回新 token
// - refresh_token 模式：消耗旧的 refresh_token，返回新 token
func (s *IQQDServiceImpl) IssueToken(ctx *gin.Context, appKey, appSecret, code, grantType, oldRefreshToken string) (models.TokenResponse, error) {
	if s.initConfig(ctx) != nil {
		return models.TokenResponse{}, models.ErrDatabaseRequired
	}
	if !s.validCredential(appKey, appSecret) {
		return models.TokenResponse{}, models.ErrInvalidCredential
	}
	switch grantType {
	case "authorization_code":
		authCode, ok := s.useCode(ctx, code)
		if !ok || authCode.AppKey != appKey {
			return models.TokenResponse{}, models.ErrInvalidAuthCode
		}
	case "refresh_token":
		item, ok := s.useRefreshToken(ctx, oldRefreshToken)
		if !ok || item.AppKey != appKey {
			return models.TokenResponse{}, models.ErrInvalidRefreshToken
		}
	default:
		return models.TokenResponse{}, models.ErrInvalidGrantType
	}

	return s.issueTokenResponse(ctx, appKey)
}

// ValidAccessToken 校验 access_token 是否有效
// 从缓存读取 token 负载，比对 token 值和 appKey，并检查是否过期
func (s *IQQDServiceImpl) ValidAccessToken(ctx *gin.Context, token, appKey string) bool {
	var item tokenPayload
	if err := s.getJSON(ctx, accessTokenKey(token), &item); err != nil {
		return false
	}
	if item.Token != token || item.AppKey != appKey {
		return false
	}
	return item.ExpiresAt.IsZero() || time.Now().Before(item.ExpiresAt)
}

// ValidSign 校验请求的 MD5 签名
// 使用缓存的 AppSecret 生成期望签名，与请求中的 sign 做不区分大小写比较
func (s *IQQDServiceImpl) ValidSign(params map[string]string, body, sign string) bool {
	if sign == "" {
		return false
	}
	expected, err := GenerateMD5Sign(params, body, s.cfg.AppSecret)
	if err != nil {
		zap.L().Error("generate qqd sign failed", zap.Error(err))
		return false
	}
	return strings.EqualFold(sign, expected)
}

// issueTokenResponse 生成新的 access_token 和 refresh_token，存入缓存并返回
func (s *IQQDServiceImpl) issueTokenResponse(ctx *gin.Context, appKey string) (models.TokenResponse, error) {
	token, err := randomToken()
	if err != nil {
		return models.TokenResponse{}, fmt.Errorf("generate token: %w", err)
	}
	newRefreshToken, err := randomToken()
	if err != nil {
		return models.TokenResponse{}, fmt.Errorf("generate refresh token: %w", err)
	}

	tokenTTL := parseDurationOrDefault(s.cfg.TokenTTL, 24*time.Hour)
	refreshTTL := parseDurationOrDefault(s.cfg.RefreshTokenTTL, 30*24*time.Hour)
	expireAt := time.Now().Add(tokenTTL)
	refreshExpireAt := time.Now().Add(refreshTTL)

	if err := s.setJSON(ctx, accessTokenKey(token), tokenPayload{AppKey: appKey, Token: token, ExpiresAt: expireAt}, tokenTTL); err != nil {
		return models.TokenResponse{}, fmt.Errorf("save token: %w", err)
	}
	if err := s.setJSON(ctx, refreshTokenKey(newRefreshToken), tokenPayload{AppKey: appKey, Token: newRefreshToken, ExpiresAt: refreshExpireAt}, refreshTTL); err != nil {
		return models.TokenResponse{}, fmt.Errorf("save refresh token: %w", err)
	}

	return models.TokenResponse{
		Token:           token,
		ExpireDate:      FormatQQDTime(expireAt),
		RefreshToken:    newRefreshToken,
		RefreshExpireAt: FormatQQDTime(refreshExpireAt),
		AppKey:          s.cfg.AppKey,
		AppSecret:       s.cfg.AppSecret,
		SelfMallAccount: s.cfg.Selfmallaccount,
	}, nil
}

// useCode 从缓存读取并删除一次性授权码（一次使用即失效）
// 返回授权码负载和是否有效，检查过期和 appKey 匹配
func (s *IQQDServiceImpl) useCode(ctx *gin.Context, code string) (authCodePayload, bool) {
	var item authCodePayload
	if err := s.getJSON(ctx, authCodeKey(code), &item); err != nil {
		return authCodePayload{}, false
	}
	s.cache.Del(ctx, authCodeKey(code))
	if !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) {
		return authCodePayload{}, false
	}
	return item, true
}

// useRefreshToken 从缓存读取并删除 refresh_token（一次使用即失效）
// 返回 token 负载和是否有效，校验 token 值和过期时间
func (s *IQQDServiceImpl) useRefreshToken(ctx *gin.Context, token string) (tokenPayload, bool) {
	var item tokenPayload
	if err := s.getJSON(ctx, refreshTokenKey(token), &item); err != nil {
		return tokenPayload{}, false
	}
	s.cache.Del(ctx, refreshTokenKey(token))
	if item.Token != token {
		return tokenPayload{}, false
	}
	if !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) {
		return tokenPayload{}, false
	}
	return item, true
}

// validCredential 使用恒定时间比较校验 appKey 和 appSecret，防止时序攻击
func (s *IQQDServiceImpl) validCredential(appKey, appSecret string) bool {
	if s.cfg.AppKey == "" || s.cfg.AppSecret == "" {
		zap.L().Error("integration config is not enabled")
		return false
	}
	return subtle.ConstantTimeCompare([]byte(appKey), []byte(s.cfg.AppKey)) == 1 &&
		subtle.ConstantTimeCompare([]byte(appSecret), []byte(s.cfg.AppSecret)) == 1
}

// setJSON 将对象序列化为 JSON 并存入缓存，支持 TTL 过期
func (s *IQQDServiceImpl) setJSON(ctx *gin.Context, key string, value any, ttl time.Duration) error {
	content, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.cache.Set(ctx, key, string(content), ttl)
	return nil
}

// getJSON 从缓存读取 JSON 字符串并反序列化到 value
func (s *IQQDServiceImpl) getJSON(ctx *gin.Context, key string, value any) error {
	content, err := s.cache.Get(ctx, key)
	if err != nil {
		if isCacheNilError(err) {
			return err
		}
		return err
	}
	return json.Unmarshal([]byte(content), value)
}

// isCacheNilError 判断是否为缓存键不存在的错误
func isCacheNilError(err error) bool {
	return err != nil && strings.EqualFold(err.Error(), "cache: nil")
}

// authCodeKey 生成授权码的缓存键
// 格式: erp_api:qqd:auth_code:<code>
func authCodeKey(code string) string {
	return cacheKeyPrefix + "auth_code:" + code
}

// accessTokenKey 生成 access_token 的缓存键
// 格式: erp_api:qqd:access_token:<token>
func accessTokenKey(token string) string {
	return cacheKeyPrefix + "access_token:" + token
}

// refreshTokenKey 生成 refresh_token 的缓存键
// 格式: erp_api:qqd:refresh_token:<token>
func refreshTokenKey(token string) string {
	return cacheKeyPrefix + "refresh_token:" + token
}

// appendCallbackParams 向 URL 追加 code 和 state 查询参数
func appendCallbackParams(rawURL, code, state string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	query.Set("code", code)
	query.Set("state", state)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

// randomToken 生成 24 字节随机数的 hex 编码字符串，用作 token 或授权码
func randomToken() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

// parseDurationOrDefault 解析 duration 字符串，失败时返回默认值
func parseDurationOrDefault(value string, fallback time.Duration) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}

// FormatQQDTime 将 time.Time 格式化为管家婆要求的日期字符串格式
// 格式: yyyy-MM-dd HH:mm:ss，零值时间返回空字符串
func FormatQQDTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// firstNonNil 返回第一个非 nil 的值
func firstNonNil(values ...any) any {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

// parseInt64Value 将任意类型的值转换为 int64，支持 int、float64、string 和 json.Number
func parseInt64Value(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case float64:
		return int64(v), nil
	case json.Number:
		return v.Int64()
	case string:
		return strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	default:
		return 0, fmt.Errorf("invalid int64 value %v", value)
	}
}

// parseStringIDValue 将任意类型的值解析为 int64 ID，等价于 parseInt64Value
func parseStringIDValue(value any) (int64, error) {
	return parseInt64Value(value)
}
