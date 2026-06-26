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

	// SyncOrder 同步订单事件。
	//
	// 实现层会把事件转为商城订单模型，并在一个事务里完成主表和子表同步。
	// 任一 DAO 操作返回错误时，service 会返回该错误并触发事务回滚。
	SyncOrder(event event.OrderEvent) error

	// SyncOrderStatus 同步订单状态
	SyncOrderStatus(event event.OrderStratusEvent) error

	// SyncOrderSend 处理订单发货通知（selfmall.order.send）。
	// 校验参数、幂等判断、拆单状态判定，事务内更新物流记录和订单状态。
	SyncOrderSend(event event.OrderSendEvent) error
}
