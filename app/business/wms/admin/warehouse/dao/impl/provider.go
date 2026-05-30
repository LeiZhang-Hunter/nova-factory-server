package impl

import "github.com/google/wire"

// ProviderSet WMS 仓储 DAO Provider。
var ProviderSet = wire.NewSet(NewWarehouseAreaDao, NewWarehouseLocationDao)
