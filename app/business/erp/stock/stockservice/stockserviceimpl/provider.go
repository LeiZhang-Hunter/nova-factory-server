package stockserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewStockService,
	NewStockCheckService,
	NewStockCheckItemService,
	NewStockInService,
	NewStockInItemService,
	NewStockMoveService,
	NewStockMoveItemService,
	NewStockOutService,
	NewStockOutItemService,
	NewStockRecordService,
)
