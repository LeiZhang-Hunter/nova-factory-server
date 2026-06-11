package guanjiapo

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"nova-factory-server/app/utils/observer/integration/api"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/observer"
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
func (o *SyncObserver) OnProductChanged(event event.ProductEvent) error {
	return nil
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
	if len(event.Orders()) == 0 {
		return errors.New("order list is empty")
	}
	order := event.Orders()[0]
	ctx := context.Background()
	_, err := o.service.OrderSyncer().SyncOrders(ctx, event)
	if err != nil {
		zap.L().Error("管家婆订单同步失败", zap.String("order", order.OrderNo()), zap.Error(err))
	}
	return err
}

//
//// OnProductChanged 商品变更同步（管家婆暂不支持）
//func (o *SyncObserver) OnProductChanged(event *event.ProductEvent) error {
//	zap.L().Debug("管家婆暂不支持商品同步", zap.String("action", string(event.Action)))
//	return nil
//}
//
//// OnStockChanged 库存变更同步（管家婆暂不支持）
//func (o *SyncObserver) OnStockChanged(event *observer.StockEvent) error {
//	zap.L().Debug("管家婆暂不支持库存同步", zap.String("action", string(event.Action)))
//	return nil
//}
