package event

// ZOrderStatusSyncReqEvent 订单状态同步请求事件，对应 emall.orderstatus.synchronize。
type ZOrderStatusSyncReqEvent interface {
	Event
	Base
	GetOrders() []ZOrderStatusSyncReqData
	GetUserId() int64
}

// ZOrderStatusSyncReqData 单条订单状态同步数据，对应 {tid, status, refundstatus}。
type ZOrderStatusSyncReqData interface {
	GetDBID() int64
	GetTid() string
	GetStatus() string
	GetRefundstatus() string
}
