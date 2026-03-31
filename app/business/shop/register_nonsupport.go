//go:build !shop

package shop

import (
	"nova-factory-server/app/business/shop/shopcontroller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	controller *shopcontroller.Controller) *Shop {
	return &Shop{}
}

func NewController() *shopcontroller.Controller {
	return &shopcontroller.Controller{}
}

var GinProviderSet = wire.NewSet(NewController, NewGinEngine)

var ProviderSet = wire.NewSet(GinProviderSet)
