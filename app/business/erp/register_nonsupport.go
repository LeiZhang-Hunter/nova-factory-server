//go:build !erp

package erp

import (
	"nova-factory-server/app/business/erp/setting/settingController"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

func NewGinEngine(app *routes.App, cache cache.Cache, controller *settingController.Controller) *Erp {
	return &Erp{}
}

func NewController() *settingController.Controller {
	return &settingController.Controller{}
}

var GinProviderSet = wire.NewSet(NewController, NewGinEngine)

var ProviderSet = wire.NewSet(GinProviderSet)
