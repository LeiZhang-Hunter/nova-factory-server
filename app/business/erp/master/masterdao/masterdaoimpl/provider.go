package masterdaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAccountDao,
	NewCustomerDao,
	NewProductDao,
	NewProductVectorDao,
	NewProductCategoryDao,
	NewProductUnitDao,
	NewSupplierDao,
	NewWarehouseDao,
)
