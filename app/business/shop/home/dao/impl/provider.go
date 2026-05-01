package impl

import "github.com/google/wire"

// ProviderSet 首页模块 DAO Provider。
var ProviderSet = wire.NewSet(NewShopHomeModuleDao, NewShopHomeModuleItemDao)
