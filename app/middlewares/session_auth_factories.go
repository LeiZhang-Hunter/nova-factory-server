package middlewares

import "nova-factory-server/app/datasource/cache"

func NewAdminAuth(cache cache.Cache) *SessionAuthBuilder {
	return NewSessionAuth(cache, AdminDomain)
}

func NewShopAuth(cache cache.Cache) *SessionAuthBuilder {
	return NewSessionAuth(cache, ShopDomain)
}

func NewOptionalShopAuth(cache cache.Cache) *SessionAuthBuilder {
	return NewShopAuth(cache).Optional()
}

func NewShopWebSocketAuth(cache cache.Cache) *SessionAuthBuilder {
	return NewShopAuth(cache).ForWebSocket()
}

func NewSessionAuthMiddlewareBuilder(cache cache.Cache) *SessionAuthBuilder {
	return NewAdminAuth(cache).WithRefresh()
}

func NewShopSessionAuthMiddlewareBuilder(cache cache.Cache) *SessionAuthBuilder {
	return NewShopAuth(cache).WithRefresh()
}

func NewShopSessionAppAuthMiddlewareBuilder(cache cache.Cache) *SessionAuthBuilder {
	return NewShopAuth(cache)
}

func NewShopSessionAppWsAuthMiddlewareBuilder(cache cache.Cache) *SessionAuthBuilder {
	return NewShopWebSocketAuth(cache)
}

func NewOptionalShopSessionAuthMiddlewareBuilder(cache cache.Cache) *SessionAuthBuilder {
	return NewOptionalShopAuth(cache)
}
