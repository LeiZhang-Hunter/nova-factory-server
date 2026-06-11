package guanjiapo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
)

// GetLoginTokenFromCache 从缓存读取管家婆登录态
func GetLoginTokenFromCache(ctx context.Context, cacheStore cache.Cache) (*OAuthTokenResponse, error) {
	if cacheStore == nil {
		return nil, errors.New("cache不能为空")
	}
	cacheKey := fmt.Sprintf(redis.IntegrationLoginCacheKeyPattern, KindGuanJiaPo)
	cacheValue, err := cacheStore.Get(ctx, cacheKey)
	if err != nil || strings.TrimSpace(cacheValue) == "" {
		return nil, err
	}
	ret := new(OAuthTokenResponse)
	if err = json.Unmarshal([]byte(cacheValue), ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SaveLoginTokenToCache 写入管家婆登录态到缓存
func SaveLoginTokenToCache(ctx context.Context, cacheStore cache.Cache, token *OAuthTokenResponse, expiration time.Duration) error {
	if cacheStore == nil {
		return errors.New("cache不能为空")
	}
	if token == nil {
		return errors.New("token不能为空")
	}
	cacheKey := fmt.Sprintf(redis.IntegrationLoginCacheKeyPattern, KindGuanJiaPo)
	value, err := json.Marshal(token)
	if err != nil {
		return err
	}
	cacheStore.Set(ctx, cacheKey, string(value), expiration)
	return nil
}

// resolveAccessToken 解析可用token，优先配置，其次缓存
func resolveAccessToken(ctx context.Context, snapshot *ConfigSnapshot, cacheStore cache.Cache) (string, error) {
	if snapshot == nil {
		return "", errors.New("管家婆配置不能为空")
	}
	token := strings.TrimSpace(snapshot.Token)
	if token == "" {
		token = strings.TrimSpace(snapshot.AccessToken)
	}
	if token != "" {
		return token, nil
	}
	if cacheStore != nil {
		cacheToken, err := GetLoginTokenFromCache(ctx, cacheStore)
		if err == nil && cacheToken != nil && strings.TrimSpace(cacheToken.Token) != "" {
			return strings.TrimSpace(cacheToken.Token), nil
		}
	}
	return "", errors.New("未找到可用token，请先完成授权")
}
