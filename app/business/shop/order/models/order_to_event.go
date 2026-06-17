package models

import (
	"nova-factory-server/app/utils/stringUtils"
	"time"
)

// BuildShopOrderSyncEvent 从数据库字段构建event数据
func BuildShopOrderSyncEvent(order *OrderSet) *OrderSyncRequest {
	if order == nil {
		return &OrderSyncRequest{}
	}
	details := make([]*OrderSyncDetail, 0, len(order.Details))
	for _, detail := range order.Details {
		if detail == nil {
			continue
		}
		details = append(details, &OrderSyncDetail{
			OID:            detail.OID,
			Barcode:        detail.Barcode,
			EShopGoodsID:   detail.EShopGoodsID,
			OuterIID:       detail.OuterIID,
			EShopGoodsName: detail.EShopGoodsName,
			EShopSKUId:     detail.EShopSkuID,
			EShopSKUName:   detail.EShopSkuName,
			NumIID:         detail.NumIID,
			SKUId:          detail.SkuID,
			Num:            detail.Num,
			Payment:        detail.Payment,
			PicPath:        detail.PicPath,
			Weight:         detail.Weight,
			Size:           detail.Size,
			UnitID:         detail.UnitID,
			UnitQty:        detail.UnitQty,
		})
	}

	accounts := make([]*OrderSyncAccount, 0, len(order.Accounts))
	for _, account := range order.Accounts {
		if account == nil {
			continue
		}
		accounts = append(accounts, &OrderSyncAccount{
			FinanceCode: account.FinanceCode,
			Total:       account.Total,
		})
	}

	created := time.Now().Format("2006-01-02 15:04:05")
	payTime := order.PayTime

	return &OrderSyncRequest{
		Orders: []*OrderSyncOrder{
			{
				Tid:              order.Tid,
				UserId:           order.UserID,
				Weight:           order.Weight,
				Size:             order.Size,
				BuyerNick:        order.BuyerNick,
				BuyerMessage:     order.BuyerMessage,
				SellerMemo:       order.SellerMemo,
				Total:            order.Total,
				Privilege:        order.Privilege,
				PostFee:          order.PostFee,
				ReceiverName:     order.ReceiverName,
				ReceiverState:    stringUtils.FirstNonEmpty(order.ReceiverProvinceName, order.ReceiverProvince),
				ReceiverCity:     stringUtils.FirstNonEmpty(order.ReceiverCityName, order.ReceiverCity),
				ReceiverDistrict: stringUtils.FirstNonEmpty(order.ReceiverDistrictName, order.ReceiverDistrict),
				ReceiverAddress:  order.ReceiverAddress,
				ReceiverPhone:    order.ReceiverPhone,
				ReceiverMobile:   order.ReceiverMobile,
				Created:          created,
				Status:           order.Status,
				Type:             order.OrderType,
				InvoiceName:      order.InvoiceName,
				SellerFlag:       order.SellerFlag,
				PayTime:          payTime,
				LogistBTypeCode:  order.LogistBTypeCode,
				LogistBillCode:   order.LogistBillCode,
				BTypeCode:        order.BTypeCode,
				Details:          details,
				Accounts:         accounts,
			},
		},
	}
}
