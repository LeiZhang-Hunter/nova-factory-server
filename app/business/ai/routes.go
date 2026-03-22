//go:build ai

package ai

import (
	"nova-factory-server/app/business/ai/aiDataSetController"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	ai *aiDataSetController.AiDataSet) *AI {

	group := app.Engine.Group("")
	ai.Dataset.PublicRoutes(group)
	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		ai.Dataset.PrivateRoutes(group)    // 工业智能体
		ai.Prediction.PrivateRoutes(group) // 工业智能体
		ai.Exception.PrivateRoutes(group)  // 工业智能体
		ai.Control.PrivateRoutes(group)
	}
	return &AI{}
}
