package grasp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/core/integration/api"
	"nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
)

// GetLoginTokenFromCache 从缓存读取管家婆登录态
func (c *Client) GetLoginTokenFromCache(ctx context.Context, cacheStore cache.Cache) (*OAuthTokenResponse, error) {
	if cacheStore == nil {
		return nil, errors.New("cache不能为空")
	}
	cacheKey := fmt.Sprintf(redis.IntegrationLoginCacheKeyPattern, api.KindGuanJiaPo)
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
func (c *Client) SaveLoginTokenToCache(ctx context.Context, cacheStore cache.Cache, token *OAuthTokenResponse, expiration time.Duration) error {
	if cacheStore == nil {
		return errors.New("cache不能为空")
	}
	if token == nil {
		return errors.New("token不能为空")
	}
	cacheBytes, err := json.Marshal(token)
	if err != nil {
		return err
	}
	cacheKey := fmt.Sprintf(redis.IntegrationLoginCacheKeyPattern, api.KindGuanJiaPo)
	cacheStore.Set(ctx, cacheKey, string(cacheBytes), expiration)
	return nil
}
