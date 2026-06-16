package event

type OrderStratusEventData interface {
	Tid() string
	Status() string
	Refundstatus() string
}

type OrderStratusEvent interface {
	Event
	Base
	Orders() []OrderStratusEventData
}
