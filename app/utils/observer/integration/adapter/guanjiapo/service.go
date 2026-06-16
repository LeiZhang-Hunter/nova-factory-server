package guanjiapo

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/auth"
	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/btype"
	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/client"
	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/order"
	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/product"
	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/stock"
	"nova-factory-server/app/utils/observer/integration/api"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/kind"
)

// Service 管家婆全渠道集成客户端
type Service struct {
	oauthURL      string
	tokenURL      string
	mode          string
	tokenSyncer   api.TokenGetter
	orderSyncer   api.OrderSyncer
	productSyncer api.Product
	stockSyncer   api.StockSearcher
	btypeSyncer   api.BtypeSearcher
}

// New 创建管家婆集成客户端
func New() api.Service {
	var oauthURL string
	var tokenURL string
	mode := viper.GetString("mode")
	if mode == "dev" || mode == "" {
		oauthURL = "http://local.gjpqqd.cn:5929"
		tokenURL = "http://local.gjpqqd.cn:5929/Service/ERPService.asmx/EMallApi"
	} else {
		oauthURL = "https://www.gjpqqd.com"
		tokenURL = "https://openapi.gjpqqd.com/Service/ERPService.asmx/EMallApi"
	}
	return &Service{
		oauthURL:      oauthURL,
		tokenURL:      tokenURL,
		tokenSyncer:   auth.New(tokenURL, oauthURL, mode),
		orderSyncer:   order.New(tokenURL, mode),
		productSyncer: product.New(tokenURL, mode),
		stockSyncer:   stock.New(tokenURL, mode),
		btypeSyncer:   btype.New(tokenURL, mode),
		mode:          strings.ToLower(mode),
	}
}

// Kind 返回当前集成类型
func (c *Service) Kind() kind.Kind {
	return KindGuanJiaPo
}

func init() {
	_ = api.Register(KindGuanJiaPo, func() api.Service {
		service := New()
		//observerInstance := newSyncObserver(service)
		//observer.GetNotifier().Register(observerInstance)
		return service
	})

}

func (c *Service) OrderSyncer() api.OrderSyncer {
	return c.orderSyncer
}

func (c *Service) TokenGetter() api.TokenGetter {
	return c.tokenSyncer
}

// ProductSearcher 返回商品查询能力
func (c *Service) ProductSearcher() api.Product {
	return c.productSyncer
}

// StockSearcher 返回库存查询能力
func (c *Service) StockSearcher() api.StockSearcher {
	return c.stockSyncer
}

// BtypeSearcher 返回往来单位查询能力
func (c *Service) BtypeSearcher() api.BtypeSearcher {
	return c.btypeSyncer
}

// CheckLoginState 返回授权地址，前端跳转后完成OAuth授权
func (c *Service) CheckLoginState(cfg config.Config, overrideRedirectURL string) (api.LoginState, error) {
	snapshot, err := client.ParseSnapshot(cfg)
	if err != nil {
		return nil, err
	}
	oauthURL, err := c.buildOAuthURL(overrideRedirectURL, snapshot)
	if err != nil {
		return nil, err
	}
	return &loginState{
		Online:   false,
		Message:  "管家婆授权页面",
		Type:     string(c.Kind()),
		CheckURL: oauthURL,
	}, nil
}

// buildOAuthURL 生成管家婆授权跳转地址
func (c *Service) buildOAuthURL(overrideRedirectURL string, snapshot *ConfigSnapshot) (string, error) {
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
