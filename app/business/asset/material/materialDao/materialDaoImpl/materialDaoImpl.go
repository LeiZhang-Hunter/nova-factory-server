package materialDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewMaterialDaoImpl)
