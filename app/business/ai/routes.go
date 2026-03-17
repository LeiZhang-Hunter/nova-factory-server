//go:build ai

package ai

import (
	"nova-factory-server/app/business/ai/aiDataSetController"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

type Ai struct{}

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	ai *aiDataSetController.AiDataSet) *Ai {

	group := app.Engine.Group("")
	ai.Dataset.PublicRoutes(group)
	ai.Dataset.PrivateRoutes(group)    // 工业智能体
	ai.Prediction.PrivateRoutes(group) // 工业智能体
	ai.Exception.PrivateRoutes(group)  // 工业智能体
	ai.Control.PrivateRoutes(group)
	return &Ai{}
}
