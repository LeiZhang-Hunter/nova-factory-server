package service

import (
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"

	"github.com/gin-gonic/gin"
)

// IOrderService ERP订单服务接口。
type IOrderService interface {
	Set(c *gin.Context, req *models.OrderSet) (*models.Order, error)
	GetByID(c *gin.Context, id uint64) (*models.Order, error)
	List(c *gin.Context, req *models.OrderQuery) (*models.OrderListData, error)
	DeleteByIDs(c *gin.Context, ids []uint64) error
	SynchronizeSalesOrders(c *gin.Context, req *models.OrderSyncRequest) (result.OrderSyncResponse, error)

	// Sync 同步销售订单
	Sync(event event.OrderEvent)
}
