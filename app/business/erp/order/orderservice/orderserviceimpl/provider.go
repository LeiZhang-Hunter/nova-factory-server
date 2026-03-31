package orderserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewOrderService)
