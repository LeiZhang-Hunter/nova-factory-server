package guanjiapo

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

// resolveOpenAPIURL 解析管家婆 OpenAPI 请求地址。
// 优先使用配置快照中的 baseUrl（追加 /openapi），否则回退到默认 tokenURL。
func resolveOpenAPIURL(defaultURL string, snapshot *ConfigSnapshot) string {
	if snapshot != nil && strings.TrimSpace(snapshot.BaseURL) != "" {
		return strings.TrimRight(strings.TrimSpace(snapshot.BaseURL), "/") + "/openapi"
	}
	return defaultURL
}

// doSignedPost 执行一次签名后的管家婆 OpenAPI POST 请求。
// 复用订单同步链路中的通用流程：构建 query 参数 -> JSON 编码 body ->
// 生成 MD5 签名 -> POST -> 读取响应。返回原始响应字节，由调用方反序列化并校验业务状态码。
// 当传输失败或 HTTP 状态非 2xx 时返回错误（错误信息包含原始响应体）。
func doSignedPost(ctx context.Context, defaultURL string, snapshot *ConfigSnapshot, token string, method string, body any) ([]byte, error) {
	if snapshot == nil {
		return nil, errors.New("管家婆配置不能为空")
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	openapiURL := resolveOpenAPIURL(defaultURL, snapshot)
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
	sign, err := generateMD5Sign(signParams, string(payload), snapshot.Credentials.AppSecret)
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
