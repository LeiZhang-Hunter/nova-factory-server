//go:build wireinject
// +build wireinject

package main

import (
	"nova-factory-server/app/business/admin"
	"nova-factory-server/app/business/ai"
	"nova-factory-server/app/business/iot"
	"nova-factory-server/app/datasource"
	"nova-factory-server/app/routes"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func finalEngine(app *routes.App, _ *admin.Admin, _ *iot.Iot, _ *ai.AI) *gin.Engine {
	return app.Engine
}

func wireApp() (*gin.Engine, func(), error) {
	panic(wire.Build(
		routes.ProviderSet,
		iot.ProviderSet,
		admin.ProviderSet,
		datasource.ProviderSet,
		finalEngine,
	))
}
