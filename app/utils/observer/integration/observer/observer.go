package observer

import (
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
)

// Observer 观察者接口，各第三方系统实现此接口以接收业务变更事件
type Observer interface {
	// Name 返回观察者名称
	Name() kind.Kind

	// OnProductChanged 商品变更回调
	OnProductChanged(event event.ProductEvent) error

	// OnStockChanged 库存变更回调
	OnStockChanged(event event.StockEvent) error

	// OnOrderChanged 订单变更回调
	OnOrderChanged(event event.OrderEvent) error
}
