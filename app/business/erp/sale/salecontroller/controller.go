package salecontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSaleOut,
	NewSaleOutItem,
	NewSaleReturn,
	NewSaleReturnItem,
	wire.Struct(new(Controller), "*"),
)

type Controller struct {
	SaleOut        *SaleOut
	SaleOutItem    *SaleOutItem
	SaleReturn     *SaleReturn
	SaleReturnItem *SaleReturnItem
}

func (c *Controller) PrivateRoutes(router *gin.RouterGroup) {
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
