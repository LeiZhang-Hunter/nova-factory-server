package qcIpqcServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewQcIpqcServiceImpl)
