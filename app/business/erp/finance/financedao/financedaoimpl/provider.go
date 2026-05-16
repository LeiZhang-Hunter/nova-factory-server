package financedaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewFinancePaymentDao,
	NewFinancePaymentItemDao,
	NewFinanceReceiptDao,
	NewFinanceReceiptItemDao,
)
