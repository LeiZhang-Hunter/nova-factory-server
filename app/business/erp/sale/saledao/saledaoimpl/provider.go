package saledaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewSaleOutDao,
	NewSaleOutItemDao,
	NewSaleReturnDao,
	NewSaleReturnItemDao,
)
