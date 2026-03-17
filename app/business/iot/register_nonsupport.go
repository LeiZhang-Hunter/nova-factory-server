//go:build !iot

package iot

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache) *Iot {
	return &Iot{}
}

var GinProviderSet = wire.NewSet(NewGinEngine)

var ProviderSet = wire.NewSet(GinProviderSet)
