package impl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewShopCombinationDao, NewShopPinkDao, NewShopSeckillDao, NewShopSeckillActivityDao, NewShopSeckillConfigDao)
