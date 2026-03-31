//go:build shop
// +build shop

package shop

import (
	"nova-factory-server/app/business/shop/shopcontroller"
	"nova-factory-server/app/business/shop/shopdao/shopdaoimpl"
	"nova-factory-server/app/business/shop/shopservice/shopserviceimpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	shopdaoimpl.ProviderSet,
	shopserviceimpl.ProviderSet,
	shopcontroller.ProviderSet,
	GinProviderSet,
)
