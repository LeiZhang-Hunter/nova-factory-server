package saleserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewSaleOutService,
	NewSaleOutItemService,
	NewSaleReturnService,
	NewSaleReturnItemService,
)
