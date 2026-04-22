//go:build erp
// +build erp

package erp

import (
	"nova-factory-server/app/business/erp/order/ordercontroller"
	"nova-factory-server/app/business/erp/setting/settingcontroller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(app *routes.App, cache cache.Cache, setting *settingcontroller.Controller, order *ordercontroller.Controller) *Erp {
	group := app.Engine.Group("")
	{
		setting.IntegrationConfig.PublicRoutes(group)
	}
	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		setting.AgentConfig.PrivateRoutes(group)
		setting.IntegrationConfig.PrivateRoutes(group)
		setting.Logistics.PrivateRoutes(group)
		order.Order.PrivateRoutes(group)
		order.OrderAudit.PrivateRoutes(group)
	}

	// mpc
	{
		// 订单审计导入mcp
		order.OrderAudit.PrivateMcpRoutes(app.McpServer)
	}
	return &Erp{}
}
