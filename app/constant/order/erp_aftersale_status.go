package order

// ERP 售后状态回调常量（selfmall.afterorder.status.sync）
const (
	ERPAfterSaleStatusWaitAgree     string = "WaitAgree"     // 待审核
	ERPAfterSaleStatusAgree         string = "Agree"         // 审核同意
	ERPAfterSaleStatusRefuse        string = "Refuse"        // 拒绝
	ERPAfterSaleStatusInvalid       string = "Invalid"       // 作废
	ERPAfterSaleStatusUnFinish      string = "UnFinish"      // 待处理完成
	ERPAfterSaleStatusFinished      string = "Finished"      // 处理完成
	ERPAfterSaleStatusAgreeInvalid  string = "AgreeInvalid"  // 同意作废
	ERPAfterSaleStatusRefuseInvalid string = "RefuseInvalid" // 拒绝作废
)

// MapERPAfterSaleStatusToLocal 将 ERP 售后状态映射为本地售后单状态。
func MapERPAfterSaleStatusToLocal(erpStatus string) (int32, bool) {
	switch erpStatus {
	case ERPAfterSaleStatusWaitAgree:
		return AftersaleStatusPendingReview, true
	case ERPAfterSaleStatusAgree:
		return AftersaleStatusApproved, true
	case ERPAfterSaleStatusRefuse:
		return AftersaleStatusRejected, true
	case ERPAfterSaleStatusInvalid:
		return AftersaleStatusRefundClosed, true
	case ERPAfterSaleStatusUnFinish:
		return AftersaleStatusRefunding, true
	case ERPAfterSaleStatusFinished:
		return AftersaleStatusPendingRefund, true
	case ERPAfterSaleStatusAgreeInvalid:
		return AftersaleStatusRefundClosed, true
	case ERPAfterSaleStatusRefuseInvalid:
		return AftersaleStatusRejected, true
	default:
		return 0, false
	}
}

// ERPAfterSaleStatusTriggersRefund 判断 ERP 状态是否应触发支付通道退款。
func ERPAfterSaleStatusTriggersRefund(erpStatus string) bool {
	return false
}
