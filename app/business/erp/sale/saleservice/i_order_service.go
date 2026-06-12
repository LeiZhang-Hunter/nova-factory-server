package saleservice

import (
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"

	"github.com/gin-gonic/gin"
)

// IOrderService ERP订单服务接口。
type IOrderService interface {
	Set(c *gin.Context, req *salemodels.OrderSet) (*salemodels.Order, error)
	GetByID(c *gin.Context, id uint64) (*salemodels.Order, error)
	List(c *gin.Context, req *salemodels.OrderQuery) (*salemodels.OrderListData, error)
	DeleteByIDs(c *gin.Context, ids []uint64) error
	CheckLoginState(c *gin.Context, req *salemodels.CheckLoginStateReq) (*salemodels.CheckLoginStateResp, error)
	SynchronizeSalesOrders(c *gin.Context, req *salemodels.OrderSyncRequest) (result.OrderSyncResponse, error)

	// Sync 同步销售订单
	Sync(event event.OrderEvent)
}
