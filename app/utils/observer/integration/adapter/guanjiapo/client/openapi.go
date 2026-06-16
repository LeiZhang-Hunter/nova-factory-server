package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ResolveOpenAPIURL 解析管家婆 OpenAPI 请求地址。
func ResolveOpenAPIURL(defaultURL string, snapshot *ConfigSnapshot) string {
	if snapshot != nil && strings.TrimSpace(snapshot.BaseURL) != "" {
		return strings.TrimRight(strings.TrimSpace(snapshot.BaseURL), "/") + "/openapi"
	}
	return defaultURL
}

// DoSignedPost 执行一次签名后的管家婆 OpenAPI POST 请求。
func DoSignedPost(ctx context.Context, defaultURL string, snapshot *ConfigSnapshot, token string, method string, body any) ([]byte, error) {
	if snapshot == nil {
		return nil, errors.New("管家婆配置不能为空")
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	openapiURL := ResolveOpenAPIURL(defaultURL, snapshot)
	parse, err := url.Parse(openapiURL)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	signParams := map[string]string{
		"app_key":     snapshot.Credentials.AppKey,
		"v":           "1.0",
		"format":      "json",
		"sign_method": "md5",
		"method":      method,
		"timestamp":   timestamp,
		"token":       token,
	}
	sign, err := GenerateMD5Sign(signParams, string(payload), snapshot.Credentials.AppSecret)
	if err != nil {
		return nil, err
	}

	params := parse.Query()
	for k, v := range signParams {
		params.Set(k, v)
	}
	params.Set("sign", sign)
	parse.RawQuery = params.Encode()

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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return respBytes, errors.New("管家婆接口请求失败: " + strings.TrimSpace(string(respBytes)))
	}
	return respBytes, nil
}
