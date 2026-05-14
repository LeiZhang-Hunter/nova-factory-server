package order

const (
	OrderStatusPending     int32 = 0  // 待支付
	OrderStatusPaid        int32 = 1  // 已支付
	OrderStatusShipped     int32 = 2  // 已发货
	OrderStatusCompleted   int32 = 3  // 已完成
	OrderStatusCancelled   int32 = 4  // 已取消
	OrderStatusPartShipped int32 = 5  // 部分发货
	OrderStatusAftersale   int32 = 99 // 售后/退款
)
