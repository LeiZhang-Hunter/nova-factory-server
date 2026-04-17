package shopdaoimpl

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewShopAddressDao, NewShopCategoryDao, NewShopGoodsDao, NewShopSkuDao)
