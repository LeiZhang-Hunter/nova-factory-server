package logistics

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewTracking)
