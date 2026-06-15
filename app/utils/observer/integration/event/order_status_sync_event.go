package event

// ZOrderStatusSyncReqEvent 订单状态同步请求事件，对应 emall.orderstatus.synchronize。
type ZOrderStatusSyncReqEvent interface {
	Event
	Base
	ZOrderStatusSyncReqDataEvent
}

// ZOrderStatusSyncReqDataEvent 订单状态同步的数据载体。
type ZOrderStatusSyncReqDataEvent interface {
	GetOrders() *[]ZOrderStatusSyncReqData
}

// ZOrderStatusSyncReqData 单条订单状态同步数据，对应 {tid, status, refundstatus}。
type ZOrderStatusSyncReqData interface {
	GetTid() string
	GetStatus() string
	GetRefundStatus() string
}
