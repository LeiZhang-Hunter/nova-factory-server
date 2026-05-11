package shop

const (
	OrderStatusPending   int32 = 0 // 待支付
	OrderStatusPaid      int32 = 1 // 已支付
	OrderStatusShipped   int32 = 2 // 已发货
	OrderStatusCompleted int32 = 3 // 已完成
	OrderStatusCancelled int32 = 4 // 已取消
)
