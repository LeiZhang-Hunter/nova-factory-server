package resourceserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewSysResourceFileService)
