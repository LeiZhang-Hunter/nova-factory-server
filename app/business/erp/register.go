//go:build erp
// +build erp

package erp

import (
	"nova-factory-server/app/business/erp/order/ordercontroller"
	"nova-factory-server/app/business/erp/order/orderdao/orderdaoimpl"
	"nova-factory-server/app/business/erp/order/orderservice/orderserviceimpl"
	"nova-factory-server/app/business/erp/setting/settingcontroller"
	"nova-factory-server/app/business/erp/setting/settingdao/settingdaoimpl"
	"nova-factory-server/app/business/erp/setting/settingservice/settingserviceimpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	orderdaoimpl.ProviderSet,
	orderserviceimpl.ProviderSet,
	ordercontroller.ProviderSet,
	settingdaoimpl.ProviderSet,
	settingserviceimpl.ProviderSet,
	settingcontroller.ProviderSet,
	GinProviderSet,
)
