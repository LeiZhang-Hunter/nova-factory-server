package materialServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewMaterialService)
