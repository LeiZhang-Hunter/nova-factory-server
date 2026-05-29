//go:build wms && erp
// +build wms,erp

package wms

import (
	warehousecontroller "nova-factory-server/app/business/wms/admin/warehouse/controller"
	warehousedaoimpl "nova-factory-server/app/business/wms/admin/warehouse/dao/impl"
	warehouseserviceimpl "nova-factory-server/app/business/wms/admin/warehouse/service/impl"

	"github.com/google/wire"
)

// ProviderSet WMS 模块 Provider。
var ProviderSet = wire.NewSet(
	warehousedaoimpl.ProviderSet,
	warehouseserviceimpl.ProviderSet,
	warehousecontroller.ProviderSet,
	GinProviderSet,
)
