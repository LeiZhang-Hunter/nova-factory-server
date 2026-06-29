package order

// 售后单状态常量
const (
	AftersaleStatusPendingReview int32 = 0 // 待审核
	AftersaleStatusApproved      int32 = 1 // 已审核
	AftersaleStatusRefunding     int32 = 2 // 退款中
	AftersaleStatusRefundSuccess int32 = 3 // 退款成功
	AftersaleStatusRefundFailed  int32 = 4 // 退款失败
	AftersaleStatusRefundClosed  int32 = 5 // 退款关闭
	AftersaleStatusRejected      int32 = 6 // 已拒绝
)

// 售后同步状态常量
const (
	AftersaleSyncPending int32 = 0 // 待同步
	AftersaleSyncSuccess int32 = 1 // 同步成功
	AftersaleSyncFailed  int32 = 2 // 同步失败
)

// IsFinalAftersaleStatus 判断售后单是否已进入终态。
func IsFinalAftersaleStatus(status int32) bool {
	switch status {
	case AftersaleStatusRefundSuccess,
		AftersaleStatusRefundFailed,
		AftersaleStatusRefundClosed,
		AftersaleStatusRejected:
		return true
	default:
		return false
	}
}

// GetAftersaleStatusText 获取售后状态文本。
func GetAftersaleStatusText(status int32) string {
	switch status {
	case AftersaleStatusPendingReview:
		return "待审核"
	case AftersaleStatusApproved:
		return "已审核"
	case AftersaleStatusRefunding:
		return "退款中"
	case AftersaleStatusRefundSuccess:
		return "退款成功"
	case AftersaleStatusRefundFailed:
		return "退款失败"
	case AftersaleStatusRefundClosed:
		return "退款关闭"
	case AftersaleStatusRejected:
		return "已拒绝"
	default:
		return "未知"
	}
}
