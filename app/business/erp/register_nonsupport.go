//go:build !erp

package erp

import (
	"nova-factory-server/app/business/erp/setting/settingcontroller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

func NewGinEngine(app *routes.App, cache cache.Cache, controller *settingcontroller.Controller) *Erp {
	return &Erp{}
}

func NewController() *settingcontroller.Controller {
	return &settingcontroller.Controller{}
}

var GinProviderSet = wire.NewSet(NewController, NewGinEngine)

var ProviderSet = wire.NewSet(GinProviderSet)
