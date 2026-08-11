package data

import (
	"nova-factory-server/app/business/data/controller"
	"nova-factory-server/app/business/data/dao"
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
	collectorDao dao.ICollectorDAO,
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
		dc.CollectorRuleController.PrivateRoutes(group)
	}

	{
		dc.PipelineRuleController.PrivateMcpRoutes(app.McpServer)
		dc.PipelineRuleAgentController.PrivateMcpRoutes(app.McpServer)
		dc.ServiceConnectionController.PrivateMcpRoutes(app.McpServer)
		dc.CollectorController.PrivateMcpRoutes(app.McpServer)
		dc.CollectorRuleController.PrivateMcpRoutes(app.McpServer)
	}

	// 设备拉取接口：走设备认证，不经过会话鉴权
	device := r.Group("")
	device.Use(middlewares.NewDeviceAuthMiddleware(collectorDao))
	device.GET("/data/collectors/rules", dc.CollectorRuleController.PullRules)

	return &Data{}
}
