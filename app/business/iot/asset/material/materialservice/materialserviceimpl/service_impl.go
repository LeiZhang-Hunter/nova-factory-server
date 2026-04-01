package materialserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewMaterialService)
