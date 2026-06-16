package order

import "strings"

const NoCod = "NoCod"
const Cod = "Cod"

func ShopStatusToOrderStatus(status int32) string {
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

func OrderStatusToShopStatus(status string) int32 {
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
