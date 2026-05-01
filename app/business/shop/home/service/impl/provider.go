package impl

import "github.com/google/wire"

// ProviderSet 首页模块 Service Provider。
var ProviderSet = wire.NewSet(NewShopHomeModuleService, NewShopHomeModuleItemService)
