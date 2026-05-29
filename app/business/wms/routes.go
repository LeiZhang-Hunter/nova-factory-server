//go:build wms && erp
// +build wms,erp

package wms

import (
	warehousecontroller "nova-factory-server/app/business/wms/admin/warehouse/controller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

// GinProviderSet WMS 路由 Provider。
var GinProviderSet = wire.NewSet(NewGinEngine)

// NewGinEngine 注册 WMS 模块路由。
func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	warehouse *warehousecontroller.Controller,
) *Wms {
	group := app.Engine.Group("")
	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		warehouse.PrivateRoutes(group)
	}
	return &Wms{}
}
