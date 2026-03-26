package orderService

import (
	"nova-factory-server/app/business/erp/order/orderModels"

	"github.com/gin-gonic/gin"
)

type IOrderService interface {
	CheckLoginState(c *gin.Context, req *orderModels.CheckLoginStateReq) (*orderModels.CheckLoginStateResp, error)
}
