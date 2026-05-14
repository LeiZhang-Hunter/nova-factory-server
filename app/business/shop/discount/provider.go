package discount

import (
	discountdaoimpl "nova-factory-server/app/business/shop/discount/dao/impl"
	discountserviceimpl "nova-factory-server/app/business/shop/discount/service/impl"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	discountdaoimpl.ProviderSet,
	discountserviceimpl.ProviderSet,
)
