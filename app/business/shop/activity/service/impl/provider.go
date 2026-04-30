package impl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewShopCombinationService, NewShopPinkService, NewShopSeckillService, NewShopSeckillActivityService, NewShopSeckillConfigService)
