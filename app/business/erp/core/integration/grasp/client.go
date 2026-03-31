package grasp

import (
	"context"

	"nova-factory-server/app/business/erp/core/integration/api"
	"nova-factory-server/app/business/erp/setting/settingmodels"

	"github.com/spf13/viper"
)

// Client 管家婆全渠道集成客户端
type Client struct {
	oauthURL string
	tokenURL string
}

// New 创建管家婆集成客户端
func New() *Client {
	var oauthURL string
	var tokenURL string
	mode := viper.GetString("mode")
	if mode == "dev" || mode == "" {
		oauthURL = "http://local.gjpqqd.cn:5929"
		tokenURL = "http://local.gjpqqd.cn:5929/Service/ERPService.asmx/EMallApi"
	} else {
		oauthURL = "https://www.gjpqqd.com"
		tokenURL = "https://www.gjpqqd.com/Service/ERPService.asmx/EMallApi"
	}
	return &Client{
		oauthURL: oauthURL,
		tokenURL: tokenURL,
	}
}

// Kind 返回当前集成类型
func (c *Client) Kind() api.Kind {
	return api.KindGuanJiaPo
}

func init() {
	_ = api.Register(api.KindGuanJiaPo, func() api.Client {
		return New()
	})
}

// CheckLoginState 返回授权地址，前端跳转后完成OAuth授权
func (c *Client) CheckLoginState(ctx context.Context, cfg *settingmodels.IntegrationConfig, overrideURL string, overrideRedirectURL string) (*api.LoginState, error) {
	snapshot, err := ParseSnapshot(cfg)
	if err != nil {
		return nil, err
	}
	oauthURL, err := c.buildOAuthURL(overrideURL, overrideRedirectURL, snapshot)
	if err != nil {
		return nil, err
	}
	return &api.LoginState{
		Online:   false,
		Message:  "管家婆授权页面",
		Type:     cfg.Type,
		CheckURL: oauthURL,
	}, nil
}
