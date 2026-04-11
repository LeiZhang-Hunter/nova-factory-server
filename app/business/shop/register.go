//go:build shop
// +build shop

package shop

import (
	"nova-factory-server/app/business/shop/product/shopcontroller"
	"nova-factory-server/app/business/shop/product/shopdao/shopdaoimpl"
	"nova-factory-server/app/business/shop/product/shopservice/shopserviceimpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	shopdaoimpl.ProviderSet,
	shopserviceimpl.ProviderSet,
	shopcontroller.ProviderSet,
	GinProviderSet,
)
