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

	"nova-factory-server/app/business/erp/setting/settingModels"
)

// OAuthTokenResponse 管家婆oauthcode换取token的返回结果
type OAuthTokenResponse struct {
	Code       int64  `json:"code"`
	Message    string `json:"message"`
	Token      string `json:"token"`
	ExpireDate string `json:"expiredate"`
	IssueDate  string `json:"issuedate"`
	AppKey     string `json:"appkey"`
	AppSecret  string `json:"appsecret"`
}

// ExchangeTokenByOAuthCode 使用oauthcode换取访问令牌
func (c *Client) ExchangeTokenByOAuthCode(ctx context.Context, cfg *settingModels.IntegrationConfig, oauthCode string) (*OAuthTokenResponse, error) {
	snapshot, err := ParseSnapshot(cfg)
	if err != nil {
		return nil, err
	}
	code := strings.TrimSpace(oauthCode)
	if code == "" {
		return nil, errors.New("oauthcode不能为空")
	}
	appKey := strings.TrimSpace(snapshot.Credentials.AppKey)
	appSecret := strings.TrimSpace(snapshot.Credentials.AppSecret)
	if appKey == "" || appSecret == "" {
		return nil, errors.New("管家婆授权参数不完整，请配置appkey/appsecret")
	}
	tokenURL := c.tokenURL
	if strings.TrimSpace(snapshot.BaseURL) != "" {
		tokenURL = strings.TrimRight(strings.TrimSpace(snapshot.BaseURL), "/") + "/openapi"
	}
	parse, err := url.Parse(tokenURL)
	if err != nil {
		return nil, err
	}
	params := parse.Query()
	params.Set("method", "emall.token.get")
	body := map[string]any{
		"appkey":    appKey,
		"appsecret": appSecret,
		"oauthcode": code,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	parse.RawQuery = params.Encode()
	targetURL := parse.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ret := &OAuthTokenResponse{}
	if err = json.Unmarshal(respBytes, ret); err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := strings.TrimSpace(ret.Message)
		if msg == "" {
			msg = string(respBytes)
		}
		return nil, errors.New("获取token失败: " + msg)
	}
	return ret, nil
}
