package stockdaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewStockDao,
	NewStockCheckDao,
	NewStockCheckItemDao,
	NewStockInDao,
	NewStockInItemDao,
	NewStockMoveDao,
	NewStockMoveItemDao,
	NewStockOutDao,
	NewStockOutItemDao,
	NewStockRecordDao,
)
