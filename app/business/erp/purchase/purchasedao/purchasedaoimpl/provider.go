package purchasedaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewPurchaseInDao,
	NewPurchaseInItemDao,
	NewPurchaseOrderDao,
	NewPurchaseOrderItemDao,
	NewPurchaseReturnDao,
	NewPurchaseReturnItemDao,
)
