package models

import (
	"nova-factory-server/app/utils/observer/integration/event"
	"time"
)

// ToOrder 转化订单 为数据库Order
func ToOrder(event event.OrderEvent) []*Order {
	if event == nil {
		return nil
	}

	eventOrders := event.GetOrders()
	if len(eventOrders) == 0 {
		return []*Order{}
	}

	orders := make([]*Order, 0, len(eventOrders))
	for _, eventOrder := range eventOrders {
		if eventOrder == nil {
			continue
		}

		order := &Order{
			Tid:                  eventOrder.GetOrderNo(),
			Weight:               eventOrder.GetWeight(),
			Size:                 eventOrder.GetSize(),
			BuyerNick:            eventOrder.GetBuyerNick(),
			BuyerMessage:         eventOrder.GetBuyerMessage(),
			SellerMemo:           eventOrder.GetSellerMemo(),
			Total:                eventOrder.GetTotalAmount(),
			Privilege:            eventOrder.GetPrivilege(),
			PostFee:              eventOrder.GetPostFee(),
			ReceiverName:         eventOrder.GetReceiverName(),
			ReceiverProvince:     eventOrder.GetReceiverState(),
			ReceiverProvinceName: eventOrder.GetReceiverState(),
			ReceiverCity:         eventOrder.GetReceiverCity(),
			ReceiverCityName:     eventOrder.GetReceiverCity(),
			ReceiverDistrict:     eventOrder.GetReceiverDistrict(),
			ReceiverDistrictName: eventOrder.GetReceiverDistrict(),
			ReceiverAddress:      eventOrder.GetReceiverAddress(),
			ReceiverPhone:        eventOrder.GetReceiverPhone(),
			ReceiverMobile:       eventOrder.GetReceiverMobile(),
			Status:               eventOrder.GetStatus(),
			Type:                 eventOrder.GetType(),
			InvoiceName:          eventOrder.GetInvoiceName(),
			SellerFlag:           eventOrder.GetSellerFlag(),
			LogistBTypeCode:      eventOrder.GetLogIstBTypeCode(),
			LogistBillCode:       eventOrder.GetLogIstBillCode(),
			BTypeCode:            eventOrder.GetBTypeCode(),
			Details:              toOrderDetails(eventOrder.GetOrderNo(), eventOrder.GetDetails()),
			Accounts:             toOrderAccounts(eventOrder.GetOrderNo(), eventOrder.GetAccounts()),
		}

		if payTime := parseOrderTime(eventOrder.GetPayTime()); payTime != nil {
			order.PayTime = payTime
		}

		orders = append(orders, order)
	}

	return orders
}

func toOrderDetails(tid string, eventDetails []event.GoodsDetail) []*OrderDetail {
	if len(eventDetails) == 0 {
		return []*OrderDetail{}
	}

	details := make([]*OrderDetail, 0, len(eventDetails))
	for _, eventDetail := range eventDetails {
		if eventDetail == nil {
			continue
		}

		details = append(details, &OrderDetail{
			Tid:            tid,
			OID:            eventDetail.GetOid(),
			Barcode:        eventDetail.GetBarcode(),
			EShopGoodsID:   eventDetail.GetEshopGoodsId(),
			OuterIID:       eventDetail.GetOuterIid(),
			EShopGoodsName: eventDetail.GetEshopGoodsName(),
			EShopSkuID:     eventDetail.GetEshopSkuId(),
			EShopSkuName:   eventDetail.GetEshopSkuName(),
			NumIID:         eventDetail.GetNumIid(),
			SkuID:          eventDetail.GetSkuId(),
			Num:            eventDetail.GetNum(),
			Payment:        eventDetail.GetPayment(),
			PicPath:        eventDetail.GetPicPath(),
			Weight:         eventDetail.GetWeight(),
			Size:           eventDetail.GetSize(),
			UnitID:         eventDetail.GetUniTid(),
			UnitQty:        eventDetail.GetUnitQty(),
		})
	}

	return details
}

func toOrderAccounts(tid string, eventAccounts []event.Account) []*OrderAccount {
	if len(eventAccounts) == 0 {
		return []*OrderAccount{}
	}

	accounts := make([]*OrderAccount, 0, len(eventAccounts))
	for _, eventAccount := range eventAccounts {
		if eventAccount == nil {
			continue
		}

		accounts = append(accounts, &OrderAccount{
			Tid:         tid,
			FinanceCode: eventAccount.GetFinanceCode(),
			Total:       eventAccount.GetTotal(),
		})
	}

	return accounts
}

func parseOrderTime(value string) *time.Time {
	if value == "" {
		return nil
	}

	layouts := []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &parsed
		}
	}

	return nil
}
