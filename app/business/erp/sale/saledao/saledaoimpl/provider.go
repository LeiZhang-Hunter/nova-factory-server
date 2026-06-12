package saledaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewOrderDao,
	NewOrderDetailDao,
	NewOrderAccountDao,
	NewOrderAuditDao,
	NewSaleOutDao,
	NewSaleOutItemDao,
	NewSaleReturnDao,
	NewSaleReturnItemDao,
)
