package shopserviceimpl

import (
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/product/shopmodels"
	orderConstant "nova-factory-server/app/constant/order"
	"strings"
)

// buildShopOrderUpdateMap 构建订单主表更新字段。
//
// 使用 map 是为了允许零值覆盖。status 字段额外经过 shouldUpdateOrderStatus 校验，
// 避免空状态、未知状态、终态覆盖和乱序事件导致的状态回退。
func buildShopOrderUpdateMap(order *shopmodels.Order, current *shopmodels.Order) map[string]any {
	updates := map[string]any{
		"weight":                 order.Weight,
		"size":                   order.Size,
		"buyer_nick":             order.BuyerNick,
		"buyer_message":          order.BuyerMessage,
		"seller_memo":            order.SellerMemo,
		"total":                  order.Total,
		"privilege":              order.Privilege,
		"post_fee":               order.PostFee,
		"receiver_name":          order.ReceiverName,
		"receiver_province":      order.ReceiverProvince,
		"receiver_province_name": order.ReceiverProvinceName,
		"receiver_city":          order.ReceiverCity,
		"receiver_city_name":     order.ReceiverCityName,
		"receiver_district":      order.ReceiverDistrict,
		"receiver_district_name": order.ReceiverDistrictName,
		"receiver_street":        order.ReceiverStreet,
		"receiver_street_name":   order.ReceiverStreetName,
		"receiver_address":       order.ReceiverAddress,
		"receiver_phone":         order.ReceiverPhone,
		"receiver_mobile":        order.ReceiverMobile,
		"receiver_zip":           order.ReceiverZip,
		"order_type":             order.Type,
		"invoice_name":           order.InvoiceName,
		"seller_flag":            order.SellerFlag,
		"pay_time":               order.PayTime,
		"logist_b_type_code":     order.LogistBTypeCode,
		"logist_bill_code":       order.LogistBillCode,
		"b_type_code":            order.BTypeCode,
		"details_json":           order.DetailsJSON,
		"accounts_json":          order.AccountsJSON,
		"transaction_id":         order.TransactionID,
		"notify_raw":             order.NotifyRaw,
		"mch_id":                 order.MchID,
		"appid":                  order.AppID,
		"payer_openid":           order.PayerOpenid,
		"bill_code":              order.BillCode,
		"sync_message":           order.SyncMessage,
		"sync_status":            order.SyncStatus,
		"sync_time":              order.SyncTime,
		"dept_id":                order.DeptID,
		"update_by":              order.UpdateBy,
		"update_time":            order.UpdateTime,
	}

	if shouldUpdateOrderStatus(current.Status, order.Status) {
		updates["status"] = strings.TrimSpace(order.Status)
	} else {
		zap.L().Debug("商城订单同步跳过状态更新",
			zap.String("tid", order.Tid),
			zap.String("current_status", current.Status),
			zap.String("incoming_status", order.Status),
		)
	}
	return updates
}

// shouldUpdateOrderStatus 判断本次同步状态是否允许覆盖数据库当前状态。
func shouldUpdateOrderStatus(current, incoming string) bool {
	current = strings.TrimSpace(current)
	incoming = strings.TrimSpace(incoming)
	if incoming == "" {
		return false
	}

	incomingRank, incomingKnown := orderStatusRank(incoming)
	if !incomingKnown {
		return false
	}
	if current == "" {
		return true
	}
	if isFinalOrderStatus(current) {
		return false
	}
	currentRank, currentKnown := orderStatusRank(current)
	if !currentKnown {
		return true
	}
	return incomingRank >= currentRank
}

// isFinalOrderStatus 判断订单是否已经进入本地终态。
func isFinalOrderStatus(status string) bool {
	switch strings.TrimSpace(status) {
	case orderConstant.ERPStatusTradeSuccess,
		orderConstant.ERPStatusTradeClosed,
		orderConstant.ERPStatusAftersale:
		return true
	default:
		return false
	}
}

// orderStatusRank 返回订单状态推进优先级。
func orderStatusRank(status string) (int, bool) {
	switch strings.TrimSpace(status) {
	case orderConstant.ERPStatusNoPay:
		return 1, true
	case orderConstant.ERPStatusPayed:
		return 2, true
	case orderConstant.ERPStatusPartSend:
		return 3, true
	case orderConstant.ERPStatusSended:
		return 4, true
	case orderConstant.ERPStatusTradeSuccess,
		orderConstant.ERPStatusTradeClosed,
		orderConstant.ERPStatusAftersale:
		return 5, true
	default:
		return 0, false
	}
}
