package callback

import (
	shopmodels "nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"
	"nova-factory-server/app/utils/store/integration"

	"go.uber.org/zap"
)

type ShopApiCallback struct {
	isErr bool
	e     *shopmodels.OrderStatusSyncEvent
}

func NewShopOrderStatusApiCallback(e *shopmodels.OrderStatusSyncEvent) *ShopApiCallback {
	return &ShopApiCallback{e: e}

}

// OnSuccess 同步处理成功时调用
func (s *ShopApiCallback) OnSuccess(t event.Event, response result.SyncProductResponse) {

}

// OnError 同步处理失败时调用
func (s *ShopApiCallback) OnError(t event.Event, response result.SyncProductResponse, err error) {
	s.isErr = true
}

// OnFinish 同步完成触发
func (s *ShopApiCallback) OnFinish(t event.Event) {
	if s.isErr {
		return
	}
	service, _, err := integration.GetStore().GetService(t.GetCtx())
	if err != nil {
		zap.L().Error("获取集成商服务失败", zap.Error(err))
		return
	}
	if service == nil {
		return
	}
	//s.e.w
	status, err := service.OrderSyncer().SyncOrderStatus(t.GetCtx(), s.e)
	if err != nil {
		return
	}
	zap.L().Info("update status success", zap.Any("status", status))
}
