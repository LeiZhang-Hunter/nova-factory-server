package shopobserver

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewShopObserver)
