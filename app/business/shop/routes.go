//go:build shop
// +build shop

package shop

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
) *Shop {
	//r := app.Engine
	//group := r.Group("")
	//不做鉴权的

	return &Shop{}
}
