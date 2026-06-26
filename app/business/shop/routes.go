//go:build shop
// +build shop

package shop

import (
	"context"

	"github.com/google/wire"
	activityController "nova-factory-server/app/business/shop/activity/controller"
	apiActivityController "nova-factory-server/app/business/shop/api/controller/activity"
	"nova-factory-server/app/business/shop/api/controller/address"
	"nova-factory-server/app/business/shop/api/controller/agent"
	"nova-factory-server/app/business/shop/api/controller/auth"
	apiCompanyController "nova-factory-server/app/business/shop/api/controller/company"
	"nova-factory-server/app/business/shop/api/controller/favorite"
	"nova-factory-server/app/business/shop/api/controller/order"
	"nova-factory-server/app/business/shop/api/controller/product"
	shopconfigController "nova-factory-server/app/business/shop/config/controller"
	shopFinanceController "nova-factory-server/app/business/shop/finance/controller"
	homeController "nova-factory-server/app/business/shop/home/controller"
	shopobserver "nova-factory-server/app/business/shop/observer"
	shopOrderController "nova-factory-server/app/business/shop/order/controller"
	shopOrderService "nova-factory-server/app/business/shop/order/service"
	"nova-factory-server/app/business/shop/product/shopcontroller"
	userController "nova-factory-server/app/business/shop/user/controller"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/routes"
	_ "nova-factory-server/app/utils/observer/integration/adapter/guanjiapo"
	"nova-factory-server/app/utils/observer/integration/observer"
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
	agentController *agent.Controller,
	productController *product.Controller,
	orderController *order.Order,
	orderNotifyController *order.OrderNotify,
	addressController *address.Address,
	shopConfigController *shopconfigController.Controller,
	favoriteController *favorite.Favorite,
	apiActivityController *apiActivityController.Controller,
	apiCompanyController *apiCompanyController.Controller,
	shopOrderController *shopOrderController.Controller,
	orderTimeoutService shopOrderService.IOrderTimeoutService,
	shopObserver *shopobserver.ShopObserver,
	shopFinance *shopFinanceController.Controller,
) *Shop {
	group := app.Engine.Group("")

	//观察者注册
	observer.GetNotifier().Register(shopObserver)

	// 启动订单超时自动取消 Consumer（主路径，每秒轮询 Redis 延迟队列）
	go orderTimeoutService.StartConsumer(context.Background())

	//不做鉴权的（可选认证：携带有效token时自动提取用户信息，用于折扣计算等场景）
	publicGroup := group.Group("")
	publicGroup.Use(middlewares.NewOptionalShopSessionAuthMiddlewareBuilder(cache).Build())
	{
		controller.Goods.PublicRoutes(publicGroup)
		authController.Auth.PublicRoutes(publicGroup)
		agentController.Voice.PublicRoutes(publicGroup)
		productController.Category.PublicRoutes(publicGroup)
		productController.Home.PublicRoutes(publicGroup)
		productController.Product.PublicRoutes(publicGroup)
		apiCompanyController.CompanyInfo.PublicRoutes(publicGroup)
		orderNotifyController.PublicRoutes(publicGroup)
		shopOrderController.AutoCancel.PublicRoutes(publicGroup)

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
		agentController.Conversations.ConfigRoutes(appGroup)
		agentController.Conversations.PrivateRoutes(appGroup)
		agentController.ASR.PrivateRoutes(appGroup)
		agentController.Message.PrivateRoutes(appGroup)
		productController.Cart.PrivateRoutes(appGroup)
		// 小程序端活动 API
		apiActivityController.Seckill.PrivateRoutes(appGroup)
		apiActivityController.Combination.PrivateRoutes(appGroup)
		apiActivityController.Pink.PrivateRoutes(appGroup)
	}

	// WebSocket 路由组 — 握手阶段即要求 Bearer 鉴权，失败时返回标准 HTTP 状态码
	wsGroup := group.Group("/api/v1/app/shop/agent/conversations")
	wsGroup.Use(middlewares.NewShopSessionAppWsAuthMiddlewareBuilder(cache).BuildForWebSocket())
	{
		agentController.Conversations.WsChatRegister(wsGroup)
		agentController.Voice.WsRegister(wsGroup)
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
		shopConfigController.Logistics.PrivateRoutes(adminGroup)
		controller.Category.PrivateRoutes(adminGroup)
		controller.Goods.PrivateRoutes(adminGroup)
		controller.Sku.PrivateRoutes(adminGroup)
		userController.DiscountRule.PrivateRoutes(adminGroup)
		userController.User.PrivateRoutes(adminGroup)
		shopOrderController.Order.PrivateRoutes(adminGroup)
		shopFinance.Account.PrivateRoutes(adminGroup)
	}

	// 加载MCP
	productController.Product.PrivateMcpRoutes(app.McpServer)

	return &Shop{}
}
