package masterserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAccountService,
	NewCustomerService,
	NewProductService,
	NewProductCategoryService,
	NewProductUnitService,
	NewSupplierService,
	NewWarehouseService,
)
