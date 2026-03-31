package orderservice

import (
	"nova-factory-server/app/business/erp/core/integration/grasp"
	"nova-factory-server/app/business/erp/order/ordermodels"

	"github.com/gin-gonic/gin"
)

type IOrderService interface {
	CheckLoginState(c *gin.Context, req *ordermodels.CheckLoginStateReq) (*ordermodels.CheckLoginStateResp, error)
	SynchronizeSalesOrders(c *gin.Context, req *grasp.OrderSyncRequest) (*grasp.OrderSyncResponse, error)
}
