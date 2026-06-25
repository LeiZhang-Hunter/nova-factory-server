package impl

import (
	"github.com/google/wire"
)

var ServiceProviderSet = wire.NewSet(NewShopSysConfigService, NewLogisticsCompanyService)
