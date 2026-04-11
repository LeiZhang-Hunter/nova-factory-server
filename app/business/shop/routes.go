//go:build shop
// +build shop

package shop

import (
	"nova-factory-server/app/business/shop/product/shopcontroller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	controller *shopcontroller.Controller,
) *Shop {
	group := app.Engine.Group("")

	//不做鉴权的
	{
		controller.Goods.PublicRoutes(group)
	}

	group.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		controller.Category.PrivateRoutes(group)
		controller.Goods.PrivateRoutes(group)
		controller.Sku.PrivateRoutes(group)
		controller.User.PrivateRoutes(group)
	}

	return &Shop{}
}
