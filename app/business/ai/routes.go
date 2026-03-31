//go:build ai

package ai

import (
	"nova-factory-server/app/business/ai/agent/aidatasetcontroller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
	"go.uber.org/zap"
)

var GinProviderSet = wire.NewSet(NewFactoryBootstrap, NewGinEngine)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	ai *aidatasetcontroller.AiDataSet,
	bootstrap *FactoryBootstrap) *AI {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.L().Warn("async init llm factories panic", zap.Any("panic", r))
			}
		}()
		if err := bootstrap.Init(); err != nil {
			zap.L().Warn("init llm factories failed", zap.Error(err))
		}
	}()

	group := app.Engine.Group("")
	ai.Dataset.PublicRoutes(group)

	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		ai.Dataset.PrivateRoutes(group)    // 工业智能体
		ai.Prediction.PrivateRoutes(group) // 工业智能体
		ai.Exception.PrivateRoutes(group)  // 工业智能体
		ai.Control.PrivateRoutes(group)
		ai.Model.PrivateRoutes(group)
		ai.OCR.PrivateRoutes(group)
	}

	return &AI{}
}
