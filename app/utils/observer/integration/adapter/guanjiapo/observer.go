package guanjiapo

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"nova-factory-server/app/utils/observer/integration/api"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/observer"
	"nova-factory-server/app/utils/observer/integration/result"
)

// SyncObserver 管家婆同步观察者，实现 observer.Observer 接口
type SyncObserver struct {
	service api.Service
}

// NewSyncObserver 创建管家婆同步观察者
func newSyncObserver(service api.Service) observer.Observer {
	return &SyncObserver{
		service: service,
	}
}

// Name 返回观察者名称
func (o *SyncObserver) Name() kind.Kind {
	return KindGuanJiaPo
}

// OnProductChanged 商品变更回调
func (o *SyncObserver) OnProductChanged(event event.ProductEvent) (result.SyncProductResponse, error) {
	return nil, nil
}

// OnStockChanged 库存变更回调
func (o *SyncObserver) OnStockChanged(event event.StockEvent) error {
	return nil
}

// OnOrderChanged 订单变更回调
func (o *SyncObserver) OnOrderChanged(event event.OrderEvent) error {
	if event.Config == nil {
		zap.L().Debug("未传入集成配置，跳过管家婆订单同步")
		return nil
	}
	if len(event.GetOrders()) == 0 {
		return errors.New("order list is empty")
	}
	order := event.GetOrders()[0]
	ctx := context.Background()
	_, err := o.service.OrderSyncer().SyncOrders(ctx, event)
	if err != nil {
		zap.L().Error("管家婆订单同步失败", zap.String("order", order.GetOrderNo()), zap.Error(err))
	}
	return err
}

// OnOrderSendChange 订单发货变化
func (o *SyncObserver) OnOrderSendChange(sendEvent event.OrderSendEvent) error {
	return nil
}
