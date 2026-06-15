package result

// AfterSaleOrderSyncResponse 售后订单同步响应接口，对应 emall.afterorder.synchronize 返回。
type AfterSaleOrderSyncResponse interface {
	GetCode() int64
	GetMessage() string
	GetOrders() []AfterSaleOrderSyncResult

	base
}

// AfterSaleOrderSyncResult 单笔售后订单同步结果接口，含 iserror 字段。
type AfterSaleOrderSyncResult interface {
	GetIsError() bool
	GetTid() string
	GetBillCode() string
	GetMessage() string

	base
}
