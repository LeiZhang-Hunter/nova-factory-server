package api

import (
	"context"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/result"
)

// Service 集成客户端最小接口，定义身份标识与登录态检查。
// 具体同步能力通过独立的小接口（OrderSyncer、ProductSyncer、StockSyncer）声明，
// 各 ERP 客户端按需实现，调用方通过类型断言判断能力。
type Service interface {
	Kind() kind.Kind
	CheckLoginState(cfg config.Config, overrideRedirectURL string) (LoginState, error)
	OrderSyncer() OrderSyncer
	TokenGetter() TokenGetter
}

// OrderSyncer 订单同步能力接口
type OrderSyncer interface {
	SyncOrders(ctx context.Context, event event.OrderEvent) (result.OrderSyncResponse, error)
}

//// ProductSyncer 商品同步能力接口
//type ProductSyncer interface {
//	SyncProducts(ctx context.Context, event event.ProductEvent) (*ProductSyncResponse, error)
//}
//
//// StockSyncer 库存同步能力接口
//type StockSyncer interface {
//	SyncStocks(ctx context.Context, event event.StockEvent) (*StockSyncResponse, error)
//}

type TokenGetter interface {
	GetToken(ctx context.Context, cfg config.Config, oauthCode string) (result.OAuthTokenResponse, error)
}
