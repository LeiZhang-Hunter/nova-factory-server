package stockcontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"nova-factory-server/app/utils/gin_mcp"
)

var ProviderSet = wire.NewSet(
	NewStock,
	NewStockCheck,
	NewStockCheckItem,
	NewStockIn,
	NewStockInItem,
	NewStockMove,
	NewStockMoveItem,
	NewStockOut,
	NewStockOutItem,
	NewStockRecord,
	wire.Struct(new(Controller), "*"),
)

type Controller struct {
	Stock          *Stock
	StockCheck     *StockCheck
	StockCheckItem *StockCheckItem
	StockIn        *StockIn
	StockInItem    *StockInItem
	StockMove      *StockMove
	StockMoveItem  *StockMoveItem
	StockOut       *StockOut
	StockOutItem   *StockOutItem
	StockRecord    *StockRecord
}

func (c *Controller) PrivateRoutes(router *gin.RouterGroup) {
	if c.Stock != nil {
		c.Stock.PrivateRoutes(router)
	}
	if c.StockCheck != nil {
		c.StockCheck.PrivateRoutes(router)
	}
	if c.StockCheckItem != nil {
		c.StockCheckItem.PrivateRoutes(router)
	}
	if c.StockIn != nil {
		c.StockIn.PrivateRoutes(router)
	}
	if c.StockInItem != nil {
		c.StockInItem.PrivateRoutes(router)
	}
	if c.StockMove != nil {
		c.StockMove.PrivateRoutes(router)
	}
	if c.StockMoveItem != nil {
		c.StockMoveItem.PrivateRoutes(router)
	}
	if c.StockOut != nil {
		c.StockOut.PrivateRoutes(router)
	}
	if c.StockOutItem != nil {
		c.StockOutItem.PrivateRoutes(router)
	}
	if c.StockRecord != nil {
		c.StockRecord.PrivateRoutes(router)
	}
}

func (c *Controller) PrivateMcpRoutes(router *gin_mcp.GinMCP) {
	if c.Stock != nil {
		c.Stock.PrivateMcpRoutes(router)
	}
	if c.StockCheck != nil {
		c.StockCheck.PrivateMcpRoutes(router)
	}
	if c.StockCheckItem != nil {
		c.StockCheckItem.PrivateMcpRoutes(router)
	}
	if c.StockIn != nil {
		c.StockIn.PrivateMcpRoutes(router)
	}
	if c.StockInItem != nil {
		c.StockInItem.PrivateMcpRoutes(router)
	}
	if c.StockMove != nil {
		c.StockMove.PrivateMcpRoutes(router)
	}
	if c.StockMoveItem != nil {
		c.StockMoveItem.PrivateMcpRoutes(router)
	}
	if c.StockOut != nil {
		c.StockOut.PrivateMcpRoutes(router)
	}
	if c.StockOutItem != nil {
		c.StockOutItem.PrivateMcpRoutes(router)
	}
	if c.StockRecord != nil {
		c.StockRecord.PrivateMcpRoutes(router)
	}
}
