//go:build erp
// +build erp

package erp

import (
	"nova-factory-server/app/business/erp/order/orderController"
	"nova-factory-server/app/business/erp/setting/settingController"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(app *routes.App, cache cache.Cache, setting *settingController.Controller, order *orderController.Controller) *Erp {
	group := app.Engine.Group("")
	{
		setting.IntegrationConfig.PublicRoutes(group)
	}
	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		setting.AgentConfig.PrivateRoutes(group)
		setting.IntegrationConfig.PrivateRoutes(group)
		order.Order.PrivateRoutes(group)
	}
	return &Erp{}
}
