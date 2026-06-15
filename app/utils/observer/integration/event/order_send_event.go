package event

// OrderSendDetail 发货详情
type OrderSendDetail interface {
	GetSubTid() string
	GetQty() int
}

// OrderSendEvent 订单发货
type OrderSendEvent interface {
	Event
	Base
	GetTid() string
	GetCompanyCode() string
	GetIsSplit() int
	GetOutSid() string
	GetSubTid() string
	GetDetails() []OrderSendDetail
}
