//go:build shop
// +build shop

package shop

import (
	activityController "nova-factory-server/app/business/shop/activity/controller"
	"nova-factory-server/app/business/shop/api/controller/auth"
	"nova-factory-server/app/business/shop/api/controller/product"
	"nova-factory-server/app/business/shop/product/shopcontroller"
	userController "nova-factory-server/app/business/shop/user/controller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(
	app *routes.App,
	cache cache.Cache,
	activityController *activityController.Controller,
	controller *shopcontroller.Controller,
	userController *userController.Controller,
	authController *auth.Controller,
	productController *product.Controller,
) *Shop {
	group := app.Engine.Group("")

	//不做鉴权的
	{
		controller.Goods.PublicRoutes(group)
		authController.Auth.PublicRoutes(group)
		productController.Category.PublicRoutes(group)
	}

	appGroup := group.Group("")
	appGroup.Use(middlewares.NewShopSessionAuthMiddlewareBuilder(cache).Build())
	{
		authController.Auth.PrivateRoutes(appGroup)
		productController.Product.PublicRoutes(appGroup)
		productController.Category.PrivateRoutes(appGroup)
	}

	adminGroup := group.Group("")
	adminGroup.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		activityController.Combination.PrivateRoutes(adminGroup)
		activityController.Pink.PrivateRoutes(adminGroup)
		userController.Address.PrivateRoutes(adminGroup)
		userController.Cart.PrivateRoutes(adminGroup)
		controller.Category.PrivateRoutes(adminGroup)
		controller.Goods.PrivateRoutes(adminGroup)
		controller.Sku.PrivateRoutes(adminGroup)
		userController.User.PrivateRoutes(adminGroup)
	}

	return &Shop{}
}
