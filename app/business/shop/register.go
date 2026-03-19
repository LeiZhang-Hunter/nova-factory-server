//go:build shop
// +build shop

package shop

import (
	"nova-factory-server/app/business/shop/shopController"
	"nova-factory-server/app/business/shop/shopDao/shopDaoImpl"
	"nova-factory-server/app/business/shop/shopService/shopServiceImpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	shopDaoImpl.ProviderSet,
	shopServiceImpl.ProviderSet,
	shopController.ProviderSet,
	GinProviderSet,
)
