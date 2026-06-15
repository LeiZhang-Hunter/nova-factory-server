// 定义统一的事件类型枚举，用于标识不同业务动作产生的事件类别。
// 所有事件都通过 EventType 区分，通知器根据事件类型调用观察者对应方法。
package event

// EventType 事件类型，表示触发同步通知的业务动作
type EventType string

const (
	// EventProductCreated 商品新建事件，当网店或 ERP 中创建新商品时触发
	EventProductCreated EventType = "product.created"
	// EventProductUpdated 商品更新事件，当商品信息（价格、名称、规格等）发生变更时触发
	EventProductUpdated EventType = "product.updated"
	// EventProductDeleted 商品删除事件，当商品被下架或删除时触发
	EventProductDeleted EventType = "product.deleted"
	// EventStockChanged 库存变更事件，当商品库存数量发生变化时触发
	EventStockChanged EventType = "stock.changed"
	// EventOrderCreated 订单创建事件，当网店生成新订单时触发
	EventOrderCreated EventType = "order.created"
	// EventOrderStatusChanged 订单状态变更事件，当订单状态（付款、发货、完成等）发生流转时触发
	EventOrderStatusChanged EventType = "order.status_changed"
	//
	EventOrderSendChanged EventType = "order.send_changed"
)
