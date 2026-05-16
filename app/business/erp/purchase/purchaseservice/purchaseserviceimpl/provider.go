package purchaseserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewPurchaseInService,
	NewPurchaseInItemService,
	NewPurchaseOrderService,
	NewPurchaseOrderItemService,
	NewPurchaseReturnService,
	NewPurchaseReturnItemService,
)
