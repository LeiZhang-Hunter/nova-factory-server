package grasp

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// buildOAuthURL 生成管家婆授权跳转地址
func (c *Client) buildOAuthURL(overrideURL string, overrideRedirectURL string, snapshot *ConfigSnapshot) (string, error) {
	if snapshot == nil {
		return "", errors.New("管家婆配置不能为空")
	}
	base := strings.TrimSpace(c.oauthURL)
	u, err := url.Parse(base + "/EMallOauth.gspx")
	if err != nil {
		return "", err
	}
	appKey := strings.TrimSpace(snapshot.Credentials.AppKey)
	appSecret := strings.TrimSpace(snapshot.Credentials.AppSecret)
	redirectURL := strings.TrimSpace(overrideRedirectURL)
	if redirectURL == "" {
		redirectURL = strings.TrimSpace(snapshot.RedirectURL)
	}
	state := strings.TrimSpace(snapshot.State)
	if appKey == "" || appSecret == "" || redirectURL == "" {
		return "", errors.New("管家婆授权参数不完整，请配置appkey/appsecret/redirect_url")
	}
	if state == "" {
		state = fmt.Sprintf("%d", time.Now().Unix())
	}
	q := u.Query()
	q.Set("appkey", appKey)
	q.Set("appsecret", appSecret)
	q.Set("state", state)
	encodedQuery := q.Encode()
	if encodedQuery != "" {
		encodedQuery += "&"
	}
	u.RawQuery = encodedQuery + "redirect_url=" + url.QueryEscape(redirectURL)
	return u.String(), nil
}
