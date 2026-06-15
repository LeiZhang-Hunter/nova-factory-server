package result

// OrderStatusSyncResponse 订单状态同步响应接口，对应 emall.orderstatus.synchronize 返回。
type OrderStatusSyncResponse interface {
	GetCode() int64
	GetMessage() string
	GetOrders() []OrderStatusSyncResult

	base
}

// OrderStatusSyncResult 单笔订单状态同步结果接口，对应 orders[{tid, message}]。
type OrderStatusSyncResult interface {
	GetTid() string
	GetMessage() string

	base
}
