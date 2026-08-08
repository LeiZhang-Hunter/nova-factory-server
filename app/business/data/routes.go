package data

import (
	"nova-factory-server/app/business/data/controller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

type Data struct{}

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	dc *controller.DataControllers,
) *Data {
	r := app.Engine
	group := r.Group("")

	// 私有路由（需要认证）
	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		dc.PipelineRuleController.PrivateRoutes(group)
		dc.PipelineRuleAgentController.PrivateRoutes(group)
		dc.ServiceConnectionController.PrivateRoutes(group)
		dc.CollectorController.PrivateRoutes(group)
		dc.ModelConfigController.PrivateRoutes(group)
	}

	{
		dc.PipelineRuleController.PrivateMcpRoutes(app.McpServer)
		dc.PipelineRuleAgentController.PrivateMcpRoutes(app.McpServer)
		dc.ServiceConnectionController.PrivateMcpRoutes(app.McpServer)
		dc.CollectorController.PrivateMcpRoutes(app.McpServer)
		dc.ModelConfigController.PrivateMcpRoutes(app.McpServer)
	}

	return &Data{}
}
