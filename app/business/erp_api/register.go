//go:build erp
// +build erp

package erp_api

import (
	qqdcontroller "nova-factory-server/app/business/erp_api/controller/qqd"
	qqddaoimpl "nova-factory-server/app/business/erp_api/dao/impl"
	qqdserviceimpl "nova-factory-server/app/business/erp_api/service/impl/qqd"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	qqddaoimpl.ProviderSet,
	qqdserviceimpl.ProviderSet,
	qqdcontroller.ProviderSet,
	GinProviderSet,
)
