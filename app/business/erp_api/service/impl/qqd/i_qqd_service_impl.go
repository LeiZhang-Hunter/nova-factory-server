package qqd

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"nova-factory-server/app/business/erp/setting/settingdao"
	"nova-factory-server/app/datasource/cache"
	"strconv"
	"strings"
	"time"

	"nova-factory-server/app/business/erp_api/dao"
	"nova-factory-server/app/business/erp_api/models"
	qqdservice "nova-factory-server/app/business/erp_api/service/qqd"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const cacheKeyPrefix = "erp_api:qqd:"

type IQQDServiceImpl struct {
	cfg         models.QQDConfig
	cache       cache.Cache
	goodsDao    dao.IQQDGoodsDao
	goodsSkuDao dao.IQQDGoodsSkuDao
	configDao   settingdao.IIntegrationConfigDao
}

type authCodePayload struct {
	AppKey      string    `json:"app_key"`
	RedirectURI string    `json:"redirect_uri"`
	State       string    `json:"state"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type tokenPayload struct {
	AppKey    string    `json:"app_key"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewIQQDServiceImpl(configDao settingdao.IIntegrationConfigDao, cache cache.Cache, goodsDao dao.IQQDGoodsDao, goodsSkuDao dao.IQQDGoodsSkuDao) qqdservice.Service {
	return &IQQDServiceImpl{
		configDao:   configDao,
		cache:       cache,
		goodsDao:    goodsDao,
		goodsSkuDao: goodsSkuDao,
	}
}
func (s *IQQDServiceImpl) initConfig(ctx *gin.Context) error {
	enabled, err := s.configDao.GetEnabled(ctx)
	if err != nil {
		zap.L().Error("get integration config fail", zap.Error(err))
		return err
	}
	if api.Kind(enabled.Type) != api.KindGuanJiaPo {
		zap.L().Error("integration config type is not qqd")
		return err
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
func (s *IQQDServiceImpl) CreateAuthorizationCallback(ctx *gin.Context, appKey, appSecret, redirectURI, state string) (string, error) {
	if s.initConfig(ctx) != nil {
		return "", qqdservice.ErrInvalidCredential
	}
	if !s.validCredential(appKey, appSecret) {
		return "", qqdservice.ErrInvalidCredential
	}
	code, err := randomToken()
	if err != nil {
		return "", fmt.Errorf("generate auth code: %w", err)
	}
	codeTTL := parseDurationOrDefault(s.cfg.CodeTTL, 10*time.Minute)
	payload := authCodePayload{
		AppKey:      appKey,
		RedirectURI: redirectURI,
		State:       state,
		ExpiresAt:   time.Now().Add(codeTTL),
	}
	if err := s.setJSON(ctx, authCodeKey(code), payload, codeTTL); err != nil {
		return "", fmt.Errorf("save auth code: %w", err)
	}
	parsed, err := url.Parse(redirectURI)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	query.Set("code", code)
	query.Set("state", state)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func (s *IQQDServiceImpl) IssueToken(ctx *gin.Context, appKey, appSecret, code, grantType, oldRefreshToken string) (qqdservice.TokenResponse, error) {
	if s.initConfig(ctx) != nil {
		return qqdservice.TokenResponse{}, qqdservice.ErrDatabaseRequired
	}
	if !s.validCredential(appKey, appSecret) {
		return qqdservice.TokenResponse{}, qqdservice.ErrInvalidCredential
	}
	switch grantType {
	case "authorization_code":
		authCode, ok := s.useCode(ctx, code)
		if !ok || authCode.AppKey != appKey {
			return qqdservice.TokenResponse{}, qqdservice.ErrInvalidAuthCode
		}
	case "refresh_token":
		item, ok := s.useRefreshToken(ctx, oldRefreshToken)
		if !ok || item.AppKey != appKey {
			return qqdservice.TokenResponse{}, qqdservice.ErrInvalidRefreshToken
		}
	default:
		return qqdservice.TokenResponse{}, qqdservice.ErrInvalidGrantType
	}

	return s.issueTokenResponse(ctx, appKey)
}

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

func (s *IQQDServiceImpl) ProductList(ctx *gin.Context, request qqdservice.ProductListRequest) (map[string]any, error) {
	if s.goodsDao == nil || s.goodsSkuDao == nil {
		return nil, qqdservice.ErrDatabaseRequired
	}
	if request.PageNo < 1 {
		request.PageNo = 1
	}
	if request.PageSize < 1 || request.PageSize > 100 {
		request.PageSize = 10
	}

	products, total, err := s.goodsDao.List(ctx, request.PageNo, request.PageSize)
	if err != nil {
		return nil, err
	}

	goodsIDs := make([]string, 0, len(products))
	for _, product := range products {
		goodsIDs = append(goodsIDs, product.GoodsID)
	}

	productSkus, err := s.goodsSkuDao.ListByGoodsIDs(ctx, goodsIDs)
	if err != nil {
		return nil, err
	}

	skusByGoodsID := make(map[string][]models.QQDProductSkuTable, len(productSkus))
	for _, sku := range productSkus {
		skusByGoodsID[sku.GoodsID] = append(skusByGoodsID[sku.GoodsID], sku)
	}

	productInfo := make([]map[string]any, 0, len(products))
	for _, product := range products {
		skus := make([]map[string]any, 0, len(skusByGoodsID[product.GoodsID]))
		for _, sku := range skusByGoodsID[product.GoodsID] {
			skus = append(skus, map[string]any{
				"skuid":      sku.SkuID,
				"skuname":    sku.SkuName,
				"productid":  sku.GoodsID,
				"outerid":    sku.OuterID,
				"price":      sku.RetailPrice,
				"quantity":   float64(sku.Quantity),
				"created":    FormatQQDTime(sku.CreateTime),
				"modified":   FormatQQDTime(sku.UpdateTime),
				"properties": "",
			})
		}

		productInfo = append(productInfo, map[string]any{
			"cid":        product.ShopCategoryID,
			"catname":    "",
			"productid":  product.GoodsID,
			"name":       product.GoodsName,
			"outerid":    product.OuterID,
			"picpath":    product.ImageURL,
			"price":      product.RetailPrice,
			"barcodestr": "",
			"created":    FormatQQDTime(product.CreateTime),
			"desc":       product.Description,
			"modified":   FormatQQDTime(product.UpdateTime),
			"status":     "up",
			"quantity":   float64(product.Quantity),
			"skus":       skus,
		})
	}

	return map[string]any{
		"iserror":      false,
		"errormsg":     "ok",
		"totalresults": int(total),
		"productinfo":  productInfo,
	}, nil
}

func (s *IQQDServiceImpl) AddProducts(ctx *gin.Context, goodsInfos []map[string]any) ([]map[string]any, error) {
	if s.goodsDao == nil || s.goodsSkuDao == nil {
		return nil, qqdservice.ErrDatabaseRequired
	}
	for _, goods := range goodsInfos {
		items := []map[string]any{goods}
		hasSkus := false
		if values, ok := goods["skus"].([]map[string]any); ok {
			items = values
			hasSkus = true
		} else if values, ok := goods["skus"].([]any); ok {
			items = make([]map[string]any, 0, len(values))
			for _, value := range values {
				if sku, ok := value.(map[string]any); ok {
					items = append(items, sku)
				}
			}
			hasSkus = len(items) > 0
		}

		productQty := int64(0)
		for _, item := range items {
			productIDValue := firstNonNil(item["skuid"], item["goodsid"])
			productID, err := parseStringIDValue(productIDValue)
			if err != nil || productID == 0 {
				return nil, fmt.Errorf("invalid product_id: %v", productIDValue)
			}
			qty, err := parseInt64Value(item["quantity"])
			if err != nil {
				return nil, fmt.Errorf("invalid stock quantity for product_id=%d: %w", productID, err)
			}
			productQty += qty
			if err := s.goodsSkuDao.UpdateQuantity(ctx, strconv.FormatInt(productID, 10), qty, hasSkus); err != nil {
				return nil, err
			}
		}

		goodsIDValue := firstNonNil(goods["goodsid"], goods["productid"])
		if goodsIDValue == nil && !hasSkus {
			goodsIDValue = firstNonNil(goods["skuid"], goods["goodsid"])
		}
		goodsID, err := parseStringIDValue(goodsIDValue)
		if err != nil || goodsID == 0 {
			return nil, fmt.Errorf("invalid goods_id: %v", goodsIDValue)
		}
		if !hasSkus {
			productQty, err = parseInt64Value(goods["quantity"])
			if err != nil {
				return nil, fmt.Errorf("invalid stock quantity for product_id=%d: %w", goodsID, err)
			}
		}
		if err := s.goodsDao.UpdateQuantity(ctx, strconv.FormatInt(goodsID, 10), productQty); err != nil {
			return nil, err
		}
	}
	return goodsInfos, nil
}

func (s *IQQDServiceImpl) ProductStockUpdate(ctx *gin.Context, request qqdservice.ProductStockUpdateRequest) (map[string]any, error) {
	if s.goodsDao == nil || s.goodsSkuDao == nil {
		return nil, qqdservice.ErrDatabaseRequired
	}
	request.ProductID = strings.TrimSpace(request.ProductID)
	request.ProductQty = strings.TrimSpace(request.ProductQty)
	request.Skus = strings.TrimSpace(request.Skus)

	if request.ProductID == "" {
		return nil, qqdservice.ErrProductIDRequired
	}

	skuUpdates := make([]models.QQDSkuStockUpdate, 0)
	productQty := int64(0)
	lockSkus := false
	if request.Skus == "" {
		if request.ProductQty == "" {
			return nil, qqdservice.ErrProductQtyRequired
		}
		qty, err := strconv.ParseInt(request.ProductQty, 10, 64)
		if err != nil {
			return nil, qqdservice.ErrInvalidProductQty
		}
		productQty = qty
		skuUpdates = append(skuUpdates, models.QQDSkuStockUpdate{SkuID: request.ProductID, Quantity: qty})
	} else {
		lockSkus = true
		for _, item := range strings.Split(request.Skus, ",") {
			parts := strings.Split(strings.TrimSpace(item), ":")
			if len(parts) != 2 {
				return nil, qqdservice.ErrInvalidSkus
			}
			skuID := strings.TrimSpace(parts[0])
			rawQty := strings.TrimSpace(parts[1])
			if skuID == "" || rawQty == "" {
				return nil, qqdservice.ErrInvalidSkus
			}
			qty, err := strconv.ParseInt(rawQty, 10, 64)
			if err != nil {
				return nil, qqdservice.ErrInvalidSkus
			}
			skuUpdates = append(skuUpdates, models.QQDSkuStockUpdate{SkuID: skuID, Quantity: qty})
			productQty += qty
		}
	}

	for _, skuUpdate := range skuUpdates {
		if _, err := strconv.ParseInt(skuUpdate.SkuID, 10, 64); err != nil {
			return nil, qqdservice.ErrInvalidSkus
		}
	}

	for _, skuUpdate := range skuUpdates {
		if err := s.goodsSkuDao.UpdateQuantity(ctx, skuUpdate.SkuID, skuUpdate.Quantity, lockSkus); err != nil {
			return nil, err
		}
	}
	if err := s.goodsDao.UpdateQuantity(ctx, request.ProductID, productQty); err != nil {
		return nil, err
	}
	return map[string]any{
		"iserror":   false,
		"errormsg":  "ok",
		"productid": request.ProductID,
	}, nil
}

func (s *IQQDServiceImpl) issueTokenResponse(ctx *gin.Context, appKey string) (qqdservice.TokenResponse, error) {
	token, err := randomToken()
	if err != nil {
		return qqdservice.TokenResponse{}, fmt.Errorf("generate token: %w", err)
	}
	newRefreshToken, err := randomToken()
	if err != nil {
		return qqdservice.TokenResponse{}, fmt.Errorf("generate refresh token: %w", err)
	}

	tokenTTL := parseDurationOrDefault(s.cfg.TokenTTL, 24*time.Hour)
	refreshTTL := parseDurationOrDefault(s.cfg.RefreshTokenTTL, 30*24*time.Hour)
	expireAt := time.Now().Add(tokenTTL)
	refreshExpireAt := time.Now().Add(refreshTTL)

	if err := s.setJSON(ctx, accessTokenKey(token), tokenPayload{AppKey: appKey, Token: token, ExpiresAt: expireAt}, tokenTTL); err != nil {
		return qqdservice.TokenResponse{}, fmt.Errorf("save token: %w", err)
	}
	if err := s.setJSON(ctx, refreshTokenKey(newRefreshToken), tokenPayload{AppKey: appKey, Token: newRefreshToken, ExpiresAt: refreshExpireAt}, refreshTTL); err != nil {
		return qqdservice.TokenResponse{}, fmt.Errorf("save refresh token: %w", err)
	}

	return qqdservice.TokenResponse{
		Token:           token,
		ExpireDate:      FormatQQDTime(expireAt),
		RefreshToken:    newRefreshToken,
		RefreshExpireAt: FormatQQDTime(refreshExpireAt),
		AppKey:          s.cfg.AppKey,
		AppSecret:       s.cfg.AppSecret,
		SelfMallAccount: s.cfg.Selfmallaccount,
	}, nil
}

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

func (s *IQQDServiceImpl) validCredential(appKey, appSecret string) bool {
	if s.cfg.AppKey == "" || s.cfg.AppSecret == "" {
		zap.L().Error("integration config is not enabled")
		return false
	}
	return subtle.ConstantTimeCompare([]byte(appKey), []byte(s.cfg.AppKey)) == 1 &&
		subtle.ConstantTimeCompare([]byte(appSecret), []byte(s.cfg.AppSecret)) == 1
}

func (s *IQQDServiceImpl) setJSON(ctx *gin.Context, key string, value any, ttl time.Duration) error {
	content, err := json.Marshal(value)
	if err != nil {
		return err
	}
	s.cache.Set(ctx, key, string(content), ttl)
	return nil
}

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

func isCacheNilError(err error) bool {
	return err != nil && strings.EqualFold(err.Error(), "cache: nil")
}

func authCodeKey(code string) string {
	return cacheKeyPrefix + "auth_code:" + code
}

func accessTokenKey(token string) string {
	return cacheKeyPrefix + "access_token:" + token
}

func refreshTokenKey(token string) string {
	return cacheKeyPrefix + "refresh_token:" + token
}

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

func randomToken() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func parseDurationOrDefault(value string, fallback time.Duration) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}

func FormatQQDTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func firstNonNil(values ...any) any {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

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

func parseStringIDValue(value any) (int64, error) {
	return parseInt64Value(value)
}
