package result

// OrderSyncResponse 订单同步响应
type OrderSyncResponse interface {
	GetCode() int64
	GetMessage() string
	GetOrders() []OrderSyncResult
	base
}

// OrderSyncResult 单笔订单同步结果
type OrderSyncResult interface {
	GetTid() string
	GetBillCode() string
	GetMessage() string
	base
}
