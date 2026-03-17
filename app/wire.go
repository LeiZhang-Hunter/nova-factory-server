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

func wireApp() (*gin.Engine, func(), error) {
	panic(wire.Build(
		admin.ProviderSet,
		iot.ProviderSet,
		ai.ProviderSet,
		datasource.ProviderSet,
		routes.ProviderSet,
	))
}
