package qcIndexServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewQcIndexServiceImpl)
