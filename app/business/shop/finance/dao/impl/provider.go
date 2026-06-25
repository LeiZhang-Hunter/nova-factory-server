package daoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAccountDao,
)
