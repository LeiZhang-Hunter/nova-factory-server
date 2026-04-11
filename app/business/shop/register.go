//go:build shop
// +build shop

package shop

import (
	"nova-factory-server/app/business/shop/product/shopcontroller"
	"nova-factory-server/app/business/shop/product/shopdao/shopdaoimpl"
	"nova-factory-server/app/business/shop/product/shopservice/shopserviceimpl"

	userController "nova-factory-server/app/business/shop/user/controller"
	userDaoImpl "nova-factory-server/app/business/shop/user/dao/impl"
	userServiceImpl "nova-factory-server/app/business/shop/user/service/impl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	shopdaoimpl.ProviderSet,
	shopserviceimpl.ProviderSet,
	shopcontroller.ProviderSet,

	userDaoImpl.ProviderSet,
	userServiceImpl.ProviderSet,
	userController.ProviderSet,
	GinProviderSet,
)
