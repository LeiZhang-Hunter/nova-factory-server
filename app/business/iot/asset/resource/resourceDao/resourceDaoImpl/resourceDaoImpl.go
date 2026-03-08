package resourceDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewSysResourceFileDao)
