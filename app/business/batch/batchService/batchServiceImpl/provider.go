package batchServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewBatchServiceImpl)
