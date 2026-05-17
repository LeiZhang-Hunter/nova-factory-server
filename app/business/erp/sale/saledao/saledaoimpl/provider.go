package saledaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewOrderDao,
	NewOrderAuditDao,
	NewSaleOutDao,
	NewSaleOutItemDao,
	NewSaleReturnDao,
	NewSaleReturnItemDao,
)
