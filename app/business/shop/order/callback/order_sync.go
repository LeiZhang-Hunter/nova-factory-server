package callback

import (
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"
	"nova-factory-server/app/utils/store/integration"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OrderOrderSyncRequestCallback 订单回调处理
type OrderOrderSyncRequestCallback struct {
	ctx        *gin.Context
	orderEvent *models.OrderSyncRequest
	isErr      bool
}

func NewOrderSyncRequestCallback(c *gin.Context, orderEvent *models.OrderSyncRequest) *OrderOrderSyncRequestCallback {
	return &OrderOrderSyncRequestCallback{
		ctx:        c,
		orderEvent: orderEvent,
	}
}
func (s *OrderOrderSyncRequestCallback) OnSuccess(T event.Event, response result.SyncProductResponse) {
	//TODO implement me
	return
}

func (s *OrderOrderSyncRequestCallback) OnError(T event.Event, response result.SyncProductResponse, err error) {
	//TODO implement me
	s.isErr = true
	return
}

// OnFinish 同步完成触发
func (s *OrderOrderSyncRequestCallback) OnFinish(ev event.Event) {
	if s.isErr {
		return
	}
	order := s.orderEvent.ToEvent()
	getService, serviceConfig, err := integration.GetStore().GetService(s.ctx)
	if err != nil {
		zap.L().Error("获取集成商服务失败", zap.Error(err))
		return
	}
	if getService == nil {
		return
	}
	if serviceConfig == nil {
		return
	}
	s.orderEvent.WithConfig(serviceConfig)
	_, err = getService.OrderSyncer().SyncOrders(s.ctx, order)
	if err != nil {
		zap.L().Error("同步订单失败", zap.Error(err))
		return
	}

}
