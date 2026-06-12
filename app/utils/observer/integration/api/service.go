// 定义第三方集成客户端的核心服务接口。
// Service 是每个集成系统（管家婆、金蝶等）必须实现的基础接口，
// 提供身份标识、登录态检查及子能力获取（订单同步、Token 管理）。
// 能力通过独立小接口声明，调用方通过类型断言按需使用。
package api

import (
	"context"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/result"
	"time"
)

// Service 集成客户端核心接口，定义第三方系统的基础能力。
// 每个集成系统（如管家婆、金蝶）均需实现此接口。
// 具体同步能力通过独立的小接口（OrderSyncer 等）声明，各 ERP 客户端按需实现，
// 调用方通过类型断言判断是否具备某项能力。
type Service interface {
	// Kind 返回当前集成系统的类型标识（如 "gjp_v1"）
	Kind() kind.Kind
	// CheckLoginState 检查登录状态，返回 OAuth 授权跳转地址或登录信息
	CheckLoginState(cfg config.Config, overrideRedirectURL string) (LoginState, error)
	// OrderSyncer 返回订单同步能力接口，无此能力时返回 nil
	OrderSyncer() OrderSyncer
	// TokenGetter 返回 Token（OAuth 令牌）管理能力接口
	TokenGetter() TokenGetter
}

// OrderSyncer 订单同步能力接口，具备此能力的 Service 可将订单推送至第三方系统。
type OrderSyncer interface {
	// SyncOrders 将订单事件中的订单数据同步至第三方系统
	SyncOrders(ctx context.Context, event event.OrderEvent) (result.OrderSyncResponse, error)
}

//// ProductSyncer 商品同步能力接口（预留，暂未启用）
//type ProductSyncer interface {
//	SyncProducts(ctx context.Context, event event.ProductEvent) (*ProductSyncResponse, error)
//}
//
//// StockSyncer 库存同步能力接口（预留，暂未启用）
//type StockSyncer interface {
//	SyncStocks(ctx context.Context, event event.StockEvent) (*StockSyncResponse, error)
//}

// TokenGetter Token 管理接口，负责 OAuth 令牌的获取、缓存读写。
// 管家婆等 OAuth 授权流程中，需通过 oauthcode 换取 token 并存入缓存复用。
type TokenGetter interface {
	// GetTokenByCode 使用 oauthcode 向第三方系统换取访问令牌
	GetTokenByCode(ctx context.Context, cfg config.Config, oauthCode string) (result.OAuthTokenResponse, error)
	// SaveTokenToCache 将 Token 保存至缓存，避免重复授权
	SaveTokenToCache(ctx context.Context, cacheStore cache.Cache, token result.OAuthTokenResponse, expiration time.Duration) error
	// GetTokenByCache 从缓存读取 Token
	GetTokenByCache(ctx context.Context, cacheStore cache.Cache) (result.OAuthTokenResponse, error)
}
