package impl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIApiShopCartDaoImpl,
	NewShopAddressDao, NewIApiShopOrderDaoImpl, NewIApiShopOrderItemDaoImpl, NewIApiShopSysConfigDaoImpl,
	NewIApiShopWechatUserDaoImpl, NewIApiShopGoodsDaoImpl, NewIApiShopSkuDaoImpl)
