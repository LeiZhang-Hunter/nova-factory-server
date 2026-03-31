package grasp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"nova-factory-server/app/business/erp/setting/settingmodels"
	"nova-factory-server/app/datasource/cache"
)

// SynchronizeOrders 调用管家婆订单同步接口
func (c *Client) SynchronizeOrders(ctx context.Context, cfg *settingmodels.IntegrationConfig, req *OrderSyncRequest, cacheStore cache.Cache) (*OrderSyncResponse, error) {
	if req == nil || len(req.Orders) == 0 {
		return nil, errors.New("orders不能为空")
	}
	snapshot, err := ParseSnapshot(cfg)
	if err != nil {
		return nil, err
	}
	token, err := c.resolveAccessToken(ctx, snapshot, cacheStore)
	if err != nil {
		return nil, err
	}
	openapiURL := c.tokenURL
	if strings.TrimSpace(snapshot.BaseURL) != "" {
		openapiURL = strings.TrimRight(strings.TrimSpace(snapshot.BaseURL), "/") + "/openapi"
	}
	parse, err := url.Parse(openapiURL)
	if err != nil {
		return nil, err
	}
	params := parse.Query()
	params.Set("method", "emall.order.synchronize")
	parse.RawQuery = params.Encode()
	body := map[string]any{
		"token":  token,
		"orders": req.Orders,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, parse.String(), bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ret := &OrderSyncResponse{}
	if err = json.Unmarshal(respBytes, ret); err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(ret.Message)
		if msg == "" {
			msg = string(respBytes)
		}
		return nil, errors.New("订单同步失败: " + msg)
	}
	return ret, nil
}

// resolveAccessToken 解析可用token，优先配置，其次缓存
func (c *Client) resolveAccessToken(ctx context.Context, snapshot *ConfigSnapshot, cacheStore cache.Cache) (string, error) {
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
		cacheToken, err := c.GetLoginTokenFromCache(ctx, cacheStore)
		if err == nil && cacheToken != nil && strings.TrimSpace(cacheToken.Token) != "" {
			return strings.TrimSpace(cacheToken.Token), nil
		}
	}
	return "", errors.New("未找到可用token，请先完成授权")
}
