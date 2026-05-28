package impl

import "github.com/google/wire"

// ProviderSet WMS 仓储 Service Provider。
var ProviderSet = wire.NewSet(NewWarehouseAreaService, NewWarehouseLocationService)
