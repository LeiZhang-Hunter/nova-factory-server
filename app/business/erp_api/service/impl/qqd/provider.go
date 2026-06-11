package qqd

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIQQDServiceImpl)
