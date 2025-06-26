package defectServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDefectServiceImpl)
