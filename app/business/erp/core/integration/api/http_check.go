package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func CheckByHTTP(ctx context.Context, cfgType string, checkURL string, snapshot *ConfigSnapshot) (*LoginState, error) {
	if strings.TrimSpace(checkURL) == "" {
		return nil, errors.New("未配置登录态检查地址")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, checkURL, nil)
	if err != nil {
		return nil, err
	}
	token := strings.TrimSpace(snapshot.Token)
	if token == "" {
		token = strings.TrimSpace(snapshot.AccessToken)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if strings.TrimSpace(snapshot.Cookie) != "" {
		req.Header.Set("Cookie", strings.TrimSpace(snapshot.Cookie))
	}
	for k, v := range snapshot.Headers {
		if strings.TrimSpace(k) == "" || strings.TrimSpace(v) == "" {
			continue
		}
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return &LoginState{
			Online:   false,
			Message:  "登录态检查失败: " + err.Error(),
			Type:     cfgType,
			CheckURL: checkURL,
		}, nil
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	body := strings.TrimSpace(string(bodyBytes))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		msg := "登录态有效"
		if body != "" {
			msg = msg + "，响应: " + body
		}
		return &LoginState{
			Online:   true,
			Message:  msg,
			Type:     cfgType,
			CheckURL: checkURL,
			Raw:      body,
		}, nil
	}
	msg := "登录态无效"
	if body != "" {
		msg = msg + "，响应: " + body
	}
	return &LoginState{
		Online:   false,
		Message:  msg,
		Type:     cfgType,
		CheckURL: checkURL,
		Raw:      body,
	}, nil
}
