package impl

import (
	"fmt"
	"nova-factory-server/app/business/shop/api/models"
	orderConstant "nova-factory-server/app/constant/order"
	shopConstant "nova-factory-server/app/constant/shop"
	"strings"
	"time"

	shopordermodels "nova-factory-server/app/business/shop/order/models"
)

func (s *IApiShopOrderServiceImpl) buildShopOrderSyncEvent(order *shopordermodels.OrderSet) *shopordermodels.OrderSyncRequest {
	if order == nil {
		return &shopordermodels.OrderSyncRequest{}
	}
	details := make([]*shopordermodels.OrderSyncDetail, 0, len(order.Details))
	for _, detail := range order.Details {
		if detail == nil {
			continue
		}
		details = append(details, &shopordermodels.OrderSyncDetail{
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

	accounts := make([]*shopordermodels.OrderSyncAccount, 0, len(order.Accounts))
	for _, account := range order.Accounts {
		if account == nil {
			continue
		}
		accounts = append(accounts, &shopordermodels.OrderSyncAccount{
			FinanceCode: account.FinanceCode,
			Total:       account.Total,
		})
	}

	created := time.Now().Format("2006-01-02 15:04:05")
	payTime := order.PayTime

	return &shopordermodels.OrderSyncRequest{
		Orders: []*shopordermodels.OrderSyncOrder{
			{
				Tid:              order.Tid,
				Weight:           order.Weight,
				Size:             order.Size,
				BuyerNick:        order.BuyerNick,
				BuyerMessage:     order.BuyerMessage,
				SellerMemo:       order.SellerMemo,
				Total:            order.Total,
				Privilege:        order.Privilege,
				PostFee:          order.PostFee,
				ReceiverName:     order.ReceiverName,
				ReceiverState:    s.firstNonEmpty(order.ReceiverProvinceName, order.ReceiverProvince),
				ReceiverCity:     s.firstNonEmpty(order.ReceiverCityName, order.ReceiverCity),
				ReceiverDistrict: s.firstNonEmpty(order.ReceiverDistrictName, order.ReceiverDistrict),
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

// buildERPOrderSet 将商城订单请求转换为 ERP 订单保存参数。
func (s *IApiShopOrderServiceImpl) buildERPOrderSet(
	orderNo string,
	shopUser *models.User,
	address *models.ShopUserAddressApp,
	cacheData *models.OrderCacheData,
	req *models.OrderCreateReq,
) *shopordermodels.OrderSet {
	details := make([]*shopordermodels.OrderDetailSet, 0, len(cacheData.Items))
	for index, item := range cacheData.Items {
		if item == nil {
			continue
		}
		details = append(details, &shopordermodels.OrderDetailSet{
			OID:            fmt.Sprintf("%s-%d", orderNo, index+1),
			EShopGoodsID:   fmt.Sprintf("%d", item.GoodsID),
			EShopGoodsName: item.GoodsName,
			EShopSkuID:     fmt.Sprintf("%d", item.SkuID),
			EShopSkuName:   item.SkuName,
			NumIID:         item.GoodsID,
			SkuID:          item.SkuID,
			Num:            float64(item.Quantity),
			Payment:        item.TotalAmount,
			PicPath:        item.ImageURL,
		})
	}

	return &shopordermodels.OrderSet{
		Tid:                  orderNo,
		BuyerNick:            s.buildOrderBuyerNick(shopUser),
		BuyerMessage:         strings.TrimSpace(req.Remark),
		SellerMemo:           strings.TrimSpace(req.Remark),
		Total:                cacheData.GoodsAmount,
		Privilege:            cacheData.DiscountAmount,
		PostFee:              cacheData.FreightAmount,
		ReceiverName:         address.ReceiverName,
		ReceiverProvince:     address.ProvinceCode,
		ReceiverProvinceName: address.ProvinceName,
		ReceiverCity:         address.CityCode,
		ReceiverCityName:     address.CityName,
		ReceiverDistrict:     address.DistrictCode,
		ReceiverDistrictName: address.DistrictName,
		ReceiverAddress:      address.DetailAddress,
		ReceiverPhone:        address.ReceiverMobile,
		ReceiverMobile:       address.ReceiverMobile,
		Status:               s.shopStatusToERPStatus(orderConstant.OrderStatusPending),
		OrderType:            shopConstant.NoCod,
		Details:              details,
		Accounts: []*shopordermodels.OrderAccountSet{
			{
				FinanceCode: "PAY_AMOUNT",
				Total:       cacheData.PayAmount,
			},
		},
	}
}
