// 金蝶（Kingdee）集成观察者，实现 observer.Observer 接口。
// 当前为占位实现，各事件回调仅记录日志，
// 具体同步逻辑待后续对接金蝶 ERP 接口时开发。
package kingdee

import (
	"go.uber.org/zap"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/observer"
)

// init 包初始化时自动将金蝶观察者注册到全局事件分发器
func init() {
	observer.GetNotifier().Register(NewSyncObserver())
}

// SyncObserver 金蝶同步观察者，实现 observer.Observer 接口。
// 注册到全局 Notifier 后，业务事件发生时会回调对应方法。
// 当前为占位实现，仅校验配置并记录日志。
type SyncObserver struct {
}

// NewSyncObserver 创建金蝶同步观察者实例
func NewSyncObserver() observer.Observer {
	return &SyncObserver{}
}

// Name 返回观察者名称（金蝶类型标识）
func (o *SyncObserver) Name() kind.Kind {
	return KindKingdee
}

// OnProductChanged 商品变更回调，当前为占位实现。
// 若无集成配置则跳过，否则记录日志（待实现具体同步）。
func (o *SyncObserver) OnProductChanged(event event.ProductEvent) error {
	if event.Config() == nil {
		zap.L().Debug("未传入集成配置，跳过金蝶商品同步")
		return nil
	}
	zap.L().Info("金蝶商品同步待实现", zap.String("action", string(event.Action())))
	return nil
}

// OnStockChanged 库存变更回调，当前为占位实现。
func (o *SyncObserver) OnStockChanged(event event.StockEvent) error {
	if event.Config == nil {
		zap.L().Debug("未传入集成配置，跳过金蝶库存同步")
		return nil
	}
	zap.L().Info("金蝶库存同步待实现", zap.String("action", string(event.Action())))
	return nil
}

// OnOrderChanged 订单变更回调，当前为占位实现。
func (o *SyncObserver) OnOrderChanged(event event.OrderEvent) error {
	if event.Config() == nil {
		zap.L().Debug("未传入集成配置，跳过金蝶订单同步")
		return nil
	}
	return nil
}
