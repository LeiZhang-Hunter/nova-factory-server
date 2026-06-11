package kingdee

import (
	"go.uber.org/zap"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/observer"
)

func init() {
	observer.GetNotifier().Register(NewSyncObserver())
}

// SyncObserver 金蝶同步观察者，实现 observer.Observer
type SyncObserver struct {
}

// NewSyncObserver 创建金蝶同步观察者
func NewSyncObserver() observer.Observer {
	return &SyncObserver{}
}

// Name 返回观察者名称
func (o *SyncObserver) Name() kind.Kind {
	return KindKingdee
}

// OnProductChanged 商品变更回调
func (o *SyncObserver) OnProductChanged(event event.ProductEvent) error {
	if event.Config() == nil {
		zap.L().Debug("未传入集成配置，跳过金蝶商品同步")
		return nil
	}
	zap.L().Info("金蝶商品同步待实现", zap.String("action", string(event.Action())))
	return nil
}

// OnStockChanged 库存变更回调
func (o *SyncObserver) OnStockChanged(event event.StockEvent) error {
	if event.Config == nil {
		zap.L().Debug("未传入集成配置，跳过金蝶库存同步")
		return nil
	}
	zap.L().Info("金蝶库存同步待实现", zap.String("action", string(event.Action())))
	return nil
}

func (o *SyncObserver) OnOrderChanged(event event.OrderEvent) error {
	if event.Config() == nil {
		zap.L().Debug("未传入集成配置，跳过金蝶订单同步")
		return nil
	}
	return nil
}
