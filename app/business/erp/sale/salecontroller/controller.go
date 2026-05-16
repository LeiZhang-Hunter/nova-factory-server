package salecontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"nova-factory-server/app/utils/gin_mcp"
)

var ProviderSet = wire.NewSet(
	NewOrder,
	NewOrderAudit,
	NewSaleOut,
	NewSaleOutItem,
	NewSaleReturn,
	NewSaleReturnItem,
	wire.Struct(new(Controller), "*"),
)

type Controller struct {
	Order          *Order
	OrderAudit     *OrderAudit
	SaleOut        *SaleOut
	SaleOutItem    *SaleOutItem
	SaleReturn     *SaleReturn
	SaleReturnItem *SaleReturnItem
}

func (c *Controller) PrivateRoutes(router *gin.RouterGroup) {
	if c.Order != nil {
		c.Order.PrivateRoutes(router)
	}
	if c.OrderAudit != nil {
		c.OrderAudit.PrivateRoutes(router)
	}
	if c.SaleOut != nil {
		c.SaleOut.PrivateRoutes(router)
	}
	if c.SaleOutItem != nil {
		c.SaleOutItem.PrivateRoutes(router)
	}
	if c.SaleReturn != nil {
		c.SaleReturn.PrivateRoutes(router)
	}
	if c.SaleReturnItem != nil {
		c.SaleReturnItem.PrivateRoutes(router)
	}
}

func (c *Controller) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	if c.OrderAudit != nil {
		c.OrderAudit.PrivateMcpRoutes(router)
	}
}
