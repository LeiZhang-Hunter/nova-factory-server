package orderDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewOrderDao)
