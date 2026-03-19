package shopServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewShopCategoryService, NewShopGoodsService, NewShopSkuService, NewShopUserService)
