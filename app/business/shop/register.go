//go:build shop
// +build shop

package shop

import (
	"nova-factory-server/app/business/shop/activity/controller"
	activityDaoImpl "nova-factory-server/app/business/shop/activity/dao/impl"
	activityServiceImpl "nova-factory-server/app/business/shop/activity/service/impl"
	apiAddress "nova-factory-server/app/business/shop/api/controller/address"
	"nova-factory-server/app/business/shop/api/controller/auth"
	apiFavorite "nova-factory-server/app/business/shop/api/controller/favorite"
	apiOrder "nova-factory-server/app/business/shop/api/controller/order"
	apiProduct "nova-factory-server/app/business/shop/api/controller/product"
	apiDaoImpl "nova-factory-server/app/business/shop/api/dao/impl"
	apiServiceImpl "nova-factory-server/app/business/shop/api/service/impl"
	shopConfigController "nova-factory-server/app/business/shop/config/controller"
	shopConfigDaoImpl "nova-factory-server/app/business/shop/config/dao/impl"
	shopConfigServiceImpl "nova-factory-server/app/business/shop/config/service/impl"
	"nova-factory-server/app/business/shop/discount"

	homeController "nova-factory-server/app/business/shop/home/controller"
	homeDaoImpl "nova-factory-server/app/business/shop/home/dao/impl"
	homeServiceImpl "nova-factory-server/app/business/shop/home/service/impl"
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
	apiAddress.ProviderSet,
	apiProduct.ProviderSet,
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

	userDaoImpl.ProviderSet,
	userServiceImpl.ProviderSet,
	userController.ProviderSet,
	discount.ProviderSet,
	GinProviderSet,
)
