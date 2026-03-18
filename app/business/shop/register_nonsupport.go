//go:build !shop

package shop

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache) *Shop {
	return &Shop{}
}

var GinProviderSet = wire.NewSet(NewGinEngine)

var ProviderSet = wire.NewSet(GinProviderSet)
