package orderController

import (
	"nova-factory-server/app/business/erp/order/orderModels"
	"nova-factory-server/app/business/erp/order/orderService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Order struct {
	service orderService.IOrderService
}

func NewOrder(service orderService.IOrderService) *Order {
	return &Order{service: service}
}

func (o *Order) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/erp/order")
	group.GET("/check-login-state", middlewares.HasPermission("erp:order:checkLoginState"), o.CheckLoginState)
}

func (o *Order) CheckLoginState(c *gin.Context) {
	req := new(orderModels.CheckLoginStateReq)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := o.service.CheckLoginState(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}
