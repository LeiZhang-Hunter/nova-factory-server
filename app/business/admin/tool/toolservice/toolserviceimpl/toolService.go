package toolserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewGenTabletService)
