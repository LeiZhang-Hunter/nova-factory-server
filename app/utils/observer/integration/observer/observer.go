// 定义观察者模式的核心接口（Observer）与事件分发器（Notifier）。
// Observer 是各第三方集成系统（管家婆、金蝶等）实现同步逻辑的统一入口，
// Notifier 负责管理所有注册的 Observer 并在事件发生时逐一通知。
package observer

import (
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/result"
)

// Observer 观察者接口，各第三方系统（管家婆、金蝶等）实现此接口以接收业务变更事件。
// 当网店发生商品变动、库存变化、订单变更时，Notifier 会调用对应方法通知所有已注册的 Observer。
// 事务由 Notifier 统一开启，并以首参 tx 显式传入；observer 需要事务时使用该 tx，否则可忽略。
type Observer interface {
	// Name 返回观察者名称（即集成系统类型标识），用于日志区分与调试
	Name() kind.Kind

	// OnProductChanged 商品变更回调，当商品创建、更新或删除时触发
	OnProductChanged(event event.ProductEvent) (result.SyncProductResponse, error)

	// OnStockChanged 库存变更回调，当库存数量发生变化时触发
	OnStockChanged(event event.StockEvent) error

	// OnOrderChanged 订单变更回调，当订单创建或状态变更（付款、发货等）时触发
	OnOrderChanged(event event.OrderEvent) error

	// OnOrderSendChange 订单发货变化
	OnOrderSendChange(sendEvent event.OrderSendEvent) error

	// OnOrderStatusChange 订单状态更新
	OnOrderStatusChange(sendEvent event.ZOrderStatusSyncReqEvent) error
}
