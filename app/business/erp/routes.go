//go:build erp
// +build erp

package erp

import (
	"nova-factory-server/app/business/erp/setting/settingController"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(app *routes.App, cache cache.Cache, controller *settingController.Controller) *Erp {
	group := app.Engine.Group("")
	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		controller.AgentConfig.PrivateRoutes(group)
	}
	return &Erp{}
}
