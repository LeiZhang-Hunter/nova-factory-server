package purchasecontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"nova-factory-server/app/utils/gin_mcp"
)

var ProviderSet = wire.NewSet(
	NewPurchaseIn,
	NewPurchaseInItem,
	NewPurchaseOrder,
	NewPurchaseOrderItem,
	NewPurchaseReturn,
	NewPurchaseReturnItem,
	wire.Struct(new(Controller), "*"),
)

type Controller struct {
	PurchaseIn         *PurchaseIn
	PurchaseInItem     *PurchaseInItem
	PurchaseOrder      *PurchaseOrder
	PurchaseOrderItem  *PurchaseOrderItem
	PurchaseReturn     *PurchaseReturn
	PurchaseReturnItem *PurchaseReturnItem
}

func (c *Controller) PrivateRoutes(router *gin.RouterGroup) {
	if c.PurchaseIn != nil {
		c.PurchaseIn.PrivateRoutes(router)
	}
	if c.PurchaseInItem != nil {
		c.PurchaseInItem.PrivateRoutes(router)
	}
	if c.PurchaseOrder != nil {
		c.PurchaseOrder.PrivateRoutes(router)
	}
	if c.PurchaseOrderItem != nil {
		c.PurchaseOrderItem.PrivateRoutes(router)
	}
	if c.PurchaseReturn != nil {
		c.PurchaseReturn.PrivateRoutes(router)
	}
	if c.PurchaseReturnItem != nil {
		c.PurchaseReturnItem.PrivateRoutes(router)
	}
}

func (c *Controller) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	if c.PurchaseIn != nil {
		c.PurchaseIn.PrivateMcpRoutes(router)
	}
	if c.PurchaseInItem != nil {
		c.PurchaseInItem.PrivateMcpRoutes(router)
	}
	if c.PurchaseOrder != nil {
		c.PurchaseOrder.PrivateMcpRoutes(router)
	}
	if c.PurchaseOrderItem != nil {
		c.PurchaseOrderItem.PrivateMcpRoutes(router)
	}
	if c.PurchaseReturn != nil {
		c.PurchaseReturn.PrivateMcpRoutes(router)
	}
	if c.PurchaseReturnItem != nil {
		c.PurchaseReturnItem.PrivateMcpRoutes(router)
	}
}
