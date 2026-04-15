package orderservice

import (
	"nova-factory-server/app/business/erp/core/integration/grasp"
	"nova-factory-server/app/business/erp/order/ordermodels"

	"github.com/gin-gonic/gin"
)

// IOrderService ERP订单服务接口。
type IOrderService interface {
	Set(c *gin.Context, req *ordermodels.OrderSet) (*ordermodels.Order, error)
	GetByID(c *gin.Context, id uint64) (*ordermodels.Order, error)
	List(c *gin.Context, req *ordermodels.OrderQuery) (*ordermodels.OrderListData, error)
	DeleteByIDs(c *gin.Context, ids []uint64) error
	CheckLoginState(c *gin.Context, req *ordermodels.CheckLoginStateReq) (*ordermodels.CheckLoginStateResp, error)
	SynchronizeSalesOrders(c *gin.Context, req *grasp.OrderSyncRequest) (*grasp.OrderSyncResponse, error)
}
