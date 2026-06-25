package serviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAccountService,
)
