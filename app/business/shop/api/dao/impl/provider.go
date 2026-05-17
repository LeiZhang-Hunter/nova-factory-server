package impl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIApiShopCartDaoImpl,
	NewShopAddressDao, NewIApiShopSysConfigDaoImpl,
	NewIApiShopWechatUserDaoImpl, NewIApiShopGoodsDaoImpl, NewIApiShopSkuDaoImpl, NewIApiShopFavoriteDaoImpl,
	NewIApiShopSeckillDaoImpl, NewIApiShopSeckillConfigDaoImpl,
	NewIApiShopCombinationDaoImpl, NewIApiShopPinkDaoImpl)
