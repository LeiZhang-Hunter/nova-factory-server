package impl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIApiShopCartServiceImpl, NewIApiShopWechatAuthServiceImpl,
	NewIApiShopAuthServiceImpl, NewIApiShopAddressServiceImpl, NewIApiShopOrderServiceImpl, NewIApiShopGoodsServiceImpl)
