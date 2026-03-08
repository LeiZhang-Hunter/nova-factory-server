package configurationServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewConfigurationServiceImpl)
