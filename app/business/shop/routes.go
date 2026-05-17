//go:build shop
// +build shop

package shop

import (
	activityController "nova-factory-server/app/business/shop/activity/controller"
	apiActivityController "nova-factory-server/app/business/shop/api/controller/activity"
	"nova-factory-server/app/business/shop/api/controller/address"
	"nova-factory-server/app/business/shop/api/controller/auth"
	"nova-factory-server/app/business/shop/api/controller/favorite"
	"nova-factory-server/app/business/shop/api/controller/order"
	"nova-factory-server/app/business/shop/api/controller/product"
	shopconfigController "nova-factory-server/app/business/shop/config/controller"
	homeController "nova-factory-server/app/business/shop/home/controller"
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
	homeController *homeController.Controller,
	controller *shopcontroller.Controller,
	userController *userController.Controller,
	authController *auth.Controller,
	productController *product.Controller,
	orderController *order.Order,
	addressController *address.Address,
	shopConfigController *shopconfigController.Controller,
	favoriteController *favorite.Favorite,
	apiActivityController *apiActivityController.Controller,
) *Shop {
	group := app.Engine.Group("")

	//不做鉴权的（可选认证：携带有效token时自动提取用户信息，用于折扣计算等场景）
	publicGroup := group.Group("")
	publicGroup.Use(middlewares.NewOptionalShopSessionAuthMiddlewareBuilder(cache).Build())
	{
		controller.Goods.PublicRoutes(publicGroup)
		authController.Auth.PublicRoutes(publicGroup)
		productController.Category.PublicRoutes(publicGroup)
		productController.Home.PublicRoutes(publicGroup)
		productController.Product.PublicRoutes(publicGroup)
	}

	appGroup := group.Group("")
	appGroup.Use(middlewares.NewShopSessionAppAuthMiddlewareBuilder(cache).
		Build())
	{
		authController.Auth.PrivateRoutes(appGroup)
		productController.Category.PrivateRoutes(appGroup)
		productController.Product.PrivateRoutes(appGroup)
		orderController.PrivateRoutes(appGroup)
		addressController.PrivateRoutes(appGroup)
		favoriteController.PrivateRoutes(appGroup)
		productController.Cart.PrivateRoutes(appGroup)
		// 小程序端活动 API
		apiActivityController.Seckill.PrivateRoutes(appGroup)
		apiActivityController.Combination.PrivateRoutes(appGroup)
		apiActivityController.Pink.PrivateRoutes(appGroup)
	}

	adminGroup := group.Group("")
	adminGroup.Use(middlewares.NewSessionAuthMiddlewareBuilder(cache).Build())
	{
		activityController.Combination.PrivateRoutes(adminGroup)
		activityController.Pink.PrivateRoutes(adminGroup)
		activityController.Seckill.PrivateRoutes(adminGroup)
		activityController.SeckillActivity.PrivateRoutes(adminGroup)
		activityController.SeckillConfig.PrivateRoutes(adminGroup)
		homeController.HomeModule.PrivateRoutes(adminGroup)
		homeController.HomeModuleItem.PrivateRoutes(adminGroup)
		userController.Address.PrivateRoutes(adminGroup)
		userController.Cart.PrivateRoutes(adminGroup)
		shopConfigController.WechatConfig.PrivateRoutes(adminGroup)
		controller.Category.PrivateRoutes(adminGroup)
		controller.Goods.PrivateRoutes(adminGroup)
		controller.Sku.PrivateRoutes(adminGroup)
		userController.DiscountRule.PrivateRoutes(adminGroup)
		userController.User.PrivateRoutes(adminGroup)
	}

	return &Shop{}
}
