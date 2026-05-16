package financeserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewFinancePaymentService,
	NewFinancePaymentItemService,
	NewFinanceReceiptService,
	NewFinanceReceiptItemService,
)
