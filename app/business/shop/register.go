//go:build shop
// +build shop

package shop

import (
	"nova-factory-server/app/business/shop/activity/controller"
	activityDaoImpl "nova-factory-server/app/business/shop/activity/dao/impl"
	activityServiceImpl "nova-factory-server/app/business/shop/activity/service/impl"
	apiActivity "nova-factory-server/app/business/shop/api/controller/activity"
	apiAddress "nova-factory-server/app/business/shop/api/controller/address"
	apiAgent "nova-factory-server/app/business/shop/api/controller/agent"
	"nova-factory-server/app/business/shop/api/controller/auth"
	apiCompany "nova-factory-server/app/business/shop/api/controller/company"
	apiFavorite "nova-factory-server/app/business/shop/api/controller/favorite"
	apiOrder "nova-factory-server/app/business/shop/api/controller/order"
	apiProduct "nova-factory-server/app/business/shop/api/controller/product"
	apiDaoImpl "nova-factory-server/app/business/shop/api/dao/impl"
	apiServiceImpl "nova-factory-server/app/business/shop/api/service/impl"
	shopConfigController "nova-factory-server/app/business/shop/config/controller"
	shopConfigDaoImpl "nova-factory-server/app/business/shop/config/dao/impl"
	shopConfigServiceImpl "nova-factory-server/app/business/shop/config/service/impl"
	"nova-factory-server/app/business/shop/discount"

	shopFinanceController "nova-factory-server/app/business/shop/finance/controller"
	shopFinanceDaoImpl "nova-factory-server/app/business/shop/finance/dao/impl"
	shopFinanceServiceImpl "nova-factory-server/app/business/shop/finance/service/impl"

	shopobserver "nova-factory-server/app/business/shop/observer"

	shopOrderController "nova-factory-server/app/business/shop/order/controller"
	shopOrderDaoImpl "nova-factory-server/app/business/shop/order/dao/impl"
	shopOrderServiceImpl "nova-factory-server/app/business/shop/order/service/impl"

	homeController "nova-factory-server/app/business/shop/home/controller"
	homeDaoImpl "nova-factory-server/app/business/shop/home/dao/impl"
	homeServiceImpl "nova-factory-server/app/business/shop/home/service/impl"

	logisticsController "nova-factory-server/app/business/shop/logistics/controller"
	logisticsDaoImpl "nova-factory-server/app/business/shop/logistics/dao/impl"
	logisticsServiceImpl "nova-factory-server/app/business/shop/logistics/service/impl"
	"nova-factory-server/app/business/shop/product/shopcontroller"
	"nova-factory-server/app/business/shop/product/shopdao/shopdaoimpl"
	"nova-factory-server/app/business/shop/product/shopservice/shopserviceimpl"

	userController "nova-factory-server/app/business/shop/user/controller"
	userDaoImpl "nova-factory-server/app/business/shop/user/dao/impl"
	userServiceImpl "nova-factory-server/app/business/shop/user/service/impl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(

	auth.ProviderSet,
	apiFavorite.ProviderSet,
	apiOrder.ProviderSet,
	apiAgent.ProviderSet,
	apiAddress.ProviderSet,
	apiProduct.ProviderSet,
	apiActivity.ProviderSet,
	apiCompany.ProviderSet,
	apiDaoImpl.ProviderSet,
	apiServiceImpl.ProviderSet,

	shopConfigDaoImpl.ProviderSet,
	shopConfigServiceImpl.ServiceProviderSet,
	shopConfigController.ProviderSet,

	controller.ProviderSet,
	activityDaoImpl.ProviderSet,
	activityServiceImpl.ProviderSet,
	homeController.ProviderSet,
	homeDaoImpl.ProviderSet,
	homeServiceImpl.ProviderSet,

	shopdaoimpl.ProviderSet,
	shopserviceimpl.ProviderSet,
	shopcontroller.ProviderSet,

	shopOrderDaoImpl.ProviderSet,
	shopOrderServiceImpl.ProviderSet,
	shopOrderController.ProviderSet,

	userDaoImpl.ProviderSet,
	userServiceImpl.ProviderSet,
	userController.ProviderSet,
	discount.ProviderSet,

	shopFinanceController.ProviderSet,
	shopFinanceDaoImpl.ProviderSet,
	shopFinanceServiceImpl.ProviderSet,

	shopobserver.ProviderSet,

	logisticsController.ProviderSet,
	logisticsDaoImpl.ProviderSet,
	logisticsServiceImpl.ServiceProviderSet,

	GinProviderSet,
)
