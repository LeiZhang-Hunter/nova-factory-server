package shopdaoimpl

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewShopCategoryDao, NewShopGoodsDao, NewShopGoodsVectorDao,
	NewShopSkuDao, NewIShopOrderDaoImpl, NewShopOrderDetailDaoImpl, NewShopOrderAccountDaoImpl,
	NewShopOrderSendDaoImpl, NewShopOrderSendDetailDaoImpl)
