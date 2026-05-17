package financecontroller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewFinancePayment,
	NewFinancePaymentItem,
	NewFinanceReceipt,
	NewFinanceReceiptItem,
	wire.Struct(new(Controller), "*"),
)

type Controller struct {
	FinancePayment     *FinancePayment
	FinancePaymentItem *FinancePaymentItem
	FinanceReceipt     *FinanceReceipt
	FinanceReceiptItem *FinanceReceiptItem
}

func (c *Controller) PrivateRoutes(router *gin.RouterGroup) {
	if c.FinancePayment != nil {
		c.FinancePayment.PrivateRoutes(router)
	}
	if c.FinancePaymentItem != nil {
		c.FinancePaymentItem.PrivateRoutes(router)
	}
	if c.FinanceReceipt != nil {
		c.FinanceReceipt.PrivateRoutes(router)
	}
	if c.FinanceReceiptItem != nil {
		c.FinanceReceiptItem.PrivateRoutes(router)
	}
}
