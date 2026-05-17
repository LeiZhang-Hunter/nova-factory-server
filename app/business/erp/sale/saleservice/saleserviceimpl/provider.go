package saleserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewOrderService,
	NewOrderAuditService,
	NewSaleOutService,
	NewSaleOutItemService,
	NewSaleReturnService,
	NewSaleReturnItemService,
)
