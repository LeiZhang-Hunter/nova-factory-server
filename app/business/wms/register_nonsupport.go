//go:build !wms

package wms

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

// NewGinEngine 创建未启用 WMS 模块时的空实现。
func NewGinEngine(app *routes.App, cache cache.Cache) *Wms {
	return &Wms{}
}

// GinProviderSet WMS 路由 Provider。
var GinProviderSet = wire.NewSet(NewGinEngine)

// ProviderSet WMS 模块 Provider。
var ProviderSet = wire.NewSet(GinProviderSet)
