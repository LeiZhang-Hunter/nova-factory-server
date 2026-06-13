//go:build erp
// +build erp

package erp

import (
	"nova-factory-server/app/business/erp/erpobserver"
	"nova-factory-server/app/business/erp/finance/financecontroller"
	"nova-factory-server/app/business/erp/master/mastercontroller"
	"nova-factory-server/app/business/erp/purchase/purchasecontroller"
	"nova-factory-server/app/business/erp/sale/salecontroller"
	"nova-factory-server/app/business/erp/setting/settingcontroller"
	"nova-factory-server/app/business/erp/stock/stockcontroller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	master *mastercontroller.Controller,
	finance *financecontroller.Controller,
	purchase *purchasecontroller.Controller,
	sale *salecontroller.Controller,
	stock *stockcontroller.Controller,
	setting *settingcontroller.Controller,
	_ *erpobserver.ErpObserver,
) *Erp {
	group := app.Engine.Group("")
	{
		setting.IntegrationConfig.PublicRoutes(group)
	}
	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		master.PrivateRoutes(group)
		finance.PrivateRoutes(group)
		purchase.PrivateRoutes(group)
		sale.PrivateRoutes(group)
		stock.PrivateRoutes(group)
		setting.AgentConfig.PrivateRoutes(group)
		setting.IntegrationConfig.PrivateRoutes(group)
		setting.Logistics.PrivateRoutes(group)
	}

	{
		sale.PrivateMcpRoutes(app.McpServer)
	}
	return &Erp{}
}
