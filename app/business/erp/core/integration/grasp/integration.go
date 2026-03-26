package grasp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/core/integration/api"
	"nova-factory-server/app/business/erp/setting/settingModels"

	"github.com/spf13/viper"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Kind() api.Kind {
	return api.KindGuanJiaPo
}

func init() {
	_ = api.Register(api.KindGuanJiaPo, func() api.Client {
		return New()
	})
}

func (c *Client) CheckLoginState(ctx context.Context, cfg *settingModels.IntegrationConfig, overrideURL string, overrideRedirectURL string) (*api.LoginState, error) {
	snapshot, err := ParseSnapshot(cfg)
	if err != nil {
		return nil, err
	}
	oauthURL, err := c.buildOAuthURL(overrideURL, overrideRedirectURL, snapshot)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, oauthURL, nil)
	if err != nil {
		return nil, err
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
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return &api.LoginState{
			Online:   false,
			Message:  "管家婆登录态检查失败: " + err.Error(),
			Type:     cfg.Type,
			CheckURL: oauthURL,
		}, nil
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	body := strings.TrimSpace(string(bodyBytes))
	location := strings.TrimSpace(resp.Header.Get("Location"))
	lowerLocation := strings.ToLower(location)
	if strings.Contains(lowerLocation, "login") || strings.Contains(lowerLocation, "signin") {
		return &api.LoginState{
			Online:   false,
			Message:  "管家婆未登录，请先完成登录",
			Type:     cfg.Type,
			CheckURL: oauthURL,
			Raw:      location,
		}, nil
	}
	if location != "" && strings.Contains(location, "code=") {
		return &api.LoginState{
			Online:   true,
			Message:  "管家婆登录态有效，已返回授权码",
			Type:     cfg.Type,
			CheckURL: oauthURL,
			Raw:      location,
		}, nil
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		msg := "管家婆登录态有效"
		if body != "" {
			msg = msg + "，响应: " + body
		}
		return &api.LoginState{
			Online:   true,
			Message:  msg,
			Type:     cfg.Type,
			CheckURL: oauthURL,
			Raw:      location,
		}, nil
	}
	msg := "管家婆登录态无效"
	if body != "" {
		msg = msg + "，响应: " + body
	}
	return &api.LoginState{
		Online:   false,
		Message:  msg,
		Type:     cfg.Type,
		CheckURL: oauthURL,
		Raw:      location,
	}, nil
}

func (c *Client) buildOAuthURL(overrideURL string, overrideRedirectURL string, snapshot *ConfigSnapshot) (string, error) {
	if snapshot == nil {
		return "", errors.New("管家婆配置不能为空")
	}
	base := strings.TrimSpace(overrideURL)
	if base == "" {
		base = strings.TrimSpace(snapshot.CheckURL)
	}
	if base == "" {
		base = strings.TrimSpace(snapshot.LoginStatusURL)
	}
	if base == "" && strings.TrimSpace(snapshot.BaseURL) != "" {
		base = strings.TrimRight(strings.TrimSpace(snapshot.BaseURL), "/") + "/EMallOauth.gspx"
	}
	if base == "" {
		mode := viper.GetString("mode")
		if mode == "dev" || mode == "" {
			base = "http://local.gjpqqd.cn:5929/EMallOauth.gspx"
		} else {
			base = "https://www.gjpqqd.com/EMallOauth.gspx"
		}
	}
	u, err := url.Parse(base)
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
