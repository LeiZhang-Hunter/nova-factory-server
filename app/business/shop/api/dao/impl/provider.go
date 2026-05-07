package impl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIApiShopCartDaoImpl,
	NewShopAddressDao, NewIShopOrderDaoImpl, NewIShopOrderItemDaoImpl, NewShopSysConfigDao,
	NewShopWechatUserDao)
