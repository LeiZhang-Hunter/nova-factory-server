package order

import "strings"

const (
	ERPStatusNoPay        string = "NoPay"        // 未付款
	ERPStatusPayed        string = "Payed"        // 已付款（货到付款）
	ERPStatusSended       string = "Sended"       // 已发货
	ERPStatusPartSend     string = "PartSend"     // 部分发货
	ERPStatusTradeSuccess string = "TradeSuccess" // 交易成功
	ERPStatusTradeClosed  string = "TradeClosed"  // 交易关闭
	ERPStatusAftersale    string = "Aftersale"    // 售后/退款
	//	订单退款状态
	REFUNDStatusNormal        string = "Normal"        // 正常
	REFUNDStatusRefunding     string = "Refunding"     // 退款中
	REFUNDStatusRefundSuccess string = "RefundSuccess" // 退款成功
	REFUNDStatusRefundClosed  string = "RefundClosed"  // 退款关闭
)

// ShopStatusToERPStatus 将商城订单状态转换为 ERP 状态值。
func ShopStatusToERPStatus(status int32) string {
	switch status {
	case OrderStatusPending:
		return ERPStatusNoPay
	case OrderStatusPaid:
		return ERPStatusPayed
	case OrderStatusShipped:
		return ERPStatusSended
	case OrderStatusPartShipped:
		return ERPStatusPartSend
	case OrderStatusCompleted:
		return ERPStatusTradeSuccess
	case OrderStatusCancelled:
		return ERPStatusTradeClosed
	case OrderStatusAftersale:
		return ERPStatusAftersale
	default:
		return ERPStatusNoPay
	}
}

// ErpStatusToShopStatus 将 ERP 状态值转换为商城订单状态。
func ErpStatusToShopStatus(status string) int32 {
	switch strings.TrimSpace(status) {
	case ERPStatusNoPay:
		return OrderStatusPending
	case ERPStatusPayed:
		return OrderStatusPaid
	case ERPStatusSended:
		return OrderStatusShipped
	case ERPStatusPartSend:
		return OrderStatusPartShipped
	case ERPStatusTradeSuccess:
		return OrderStatusCompleted
	case ERPStatusTradeClosed:
		return OrderStatusCancelled
	case ERPStatusAftersale:
		return OrderStatusAftersale
	default:
		return OrderStatusPending
	}
}

// GetStatusText 获取状态文本
func GetStatusText(status int32) string {
	switch status {
	case OrderStatusPending:
		return "待支付"
	case OrderStatusPaid:
		return "已支付"
	case OrderStatusShipped:
		return "已发货"
	case OrderStatusPartShipped:
		return "部分发货"
	case OrderStatusCompleted:
		return "已完成"
	case OrderStatusCancelled:
		return "已取消"
	case OrderStatusAftersale:
		return "售后"
	default:
		return "未知"
	}
}

// OrderStatusRank 返回订单状态推进优先级。
func OrderStatusRank(status string) (int, bool) {
	switch strings.TrimSpace(status) {
	case ERPStatusNoPay:
		return 1, true
	case ERPStatusPayed:
		return 2, true
	case ERPStatusPartSend:
		return 3, true
	case ERPStatusSended:
		return 4, true
	case ERPStatusTradeSuccess,
		ERPStatusTradeClosed,
		ERPStatusAftersale:
		return 5, true
	default:
		return 0, false
	}
}

// IsFinalOrderStatus 判断订单是否已经进入本地终态。
func IsFinalOrderStatus(status string) bool {
	switch strings.TrimSpace(status) {
	case ERPStatusTradeSuccess,
		ERPStatusTradeClosed,
		ERPStatusAftersale:
		return true
	default:
		return false
	}
}
