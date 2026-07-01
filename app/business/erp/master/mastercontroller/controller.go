package mastercontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"nova-factory-server/app/utils/gin_mcp"
)

var ProviderSet = wire.NewSet(
	NewAccount,
	NewCustomer,
	NewProduct,
	NewProductCategory,
	NewProductUnit,
	NewSupplier,
	NewWarehouse,
	wire.Struct(new(Controller), "*"),
)

type Controller struct {
	Account         *Account
	Customer        *Customer
	Product         *Product
	ProductCategory *ProductCategory
	ProductUnit     *ProductUnit
	Supplier        *Supplier
	Warehouse       *Warehouse
}

func (c *Controller) PrivateRoutes(router *gin.RouterGroup) {
	if c.Account != nil {
		c.Account.PrivateRoutes(router)
	}
	if c.Customer != nil {
		c.Customer.PrivateRoutes(router)
	}
	if c.Product != nil {
		c.Product.PrivateRoutes(router)
	}
	if c.ProductCategory != nil {
		c.ProductCategory.PrivateRoutes(router)
	}
	if c.ProductUnit != nil {
		c.ProductUnit.PrivateRoutes(router)
	}
	if c.Supplier != nil {
		c.Supplier.PrivateRoutes(router)
	}
	if c.Warehouse != nil {
		c.Warehouse.PrivateRoutes(router)
	}
}

func (c *Controller) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	if c.Account != nil {
		c.Account.PrivateMcpRoutes(router)
	}
	if c.Customer != nil {
		c.Customer.PrivateMcpRoutes(router)
	}
	if c.Product != nil {
		c.Product.PrivateMcpRoutes(router)
	}
	if c.ProductCategory != nil {
		c.ProductCategory.PrivateMcpRoutes(router)
	}
	if c.ProductUnit != nil {
		c.ProductUnit.PrivateMcpRoutes(router)
	}
	if c.Supplier != nil {
		c.Supplier.PrivateMcpRoutes(router)
	}
	if c.Warehouse != nil {
		c.Warehouse.PrivateMcpRoutes(router)
	}
}
