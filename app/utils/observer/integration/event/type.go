package event

// EventType 事件类型
type EventType string

const (
	EventProductCreated     EventType = "product.created"
	EventProductUpdated     EventType = "product.updated"
	EventProductDeleted     EventType = "product.deleted"
	EventStockChanged       EventType = "stock.changed"
	EventOrderCreated       EventType = "order.created"
	EventOrderStatusChanged EventType = "order.status_changed"
)
