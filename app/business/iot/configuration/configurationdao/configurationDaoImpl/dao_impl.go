package configurationDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewConfigurationDaoImpl)
