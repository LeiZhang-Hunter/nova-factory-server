//go:build erp
// +build erp

package erp

import (
	"nova-factory-server/app/business/erp/order/orderController"
	"nova-factory-server/app/business/erp/order/orderDao/orderDaoImpl"
	"nova-factory-server/app/business/erp/order/orderService/orderServiceImpl"
	"nova-factory-server/app/business/erp/setting/settingController"
	"nova-factory-server/app/business/erp/setting/settingDao/settingDaoImpl"
	"nova-factory-server/app/business/erp/setting/settingService/settingServiceImpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	orderDaoImpl.ProviderSet,
	orderServiceImpl.ProviderSet,
	orderController.ProviderSet,
	settingDaoImpl.ProviderSet,
	settingServiceImpl.ProviderSet,
	settingController.ProviderSet,
	GinProviderSet,
)
