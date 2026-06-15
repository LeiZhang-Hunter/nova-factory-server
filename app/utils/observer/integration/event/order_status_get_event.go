package event

// ZOrderStatusGetReqEvent 订单状态查询请求事件，对应 emall.orderstatus.get。
type ZOrderStatusGetReqEvent struct {
	Event
	Base
	OrderStatusGetReqDataEvent
}

// OrderStatusGetReqDataEvent 订单状态查询数据载体。
type OrderStatusGetReqDataEvent interface {
	GetOrderCodes() *[]string
}
