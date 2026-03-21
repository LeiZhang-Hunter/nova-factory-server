//go:build erp
// +build erp

package erp

import (
	"nova-factory-server/app/business/erp/setting/settingController"
	"nova-factory-server/app/business/erp/setting/settingDao/settingDaoImpl"
	"nova-factory-server/app/business/erp/setting/settingService/settingServiceImpl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	settingDaoImpl.ProviderSet,
	settingServiceImpl.ProviderSet,
	settingController.ProviderSet,
	GinProviderSet,
)
