package event

type OrderStratusEventData interface {
	GetTid() string
	GetStatus() string
	GetRefundstatus() string
}

type OrderStratusEvent interface {
	Event
	Base
	GetOrders() []OrderStratusEventData
}
