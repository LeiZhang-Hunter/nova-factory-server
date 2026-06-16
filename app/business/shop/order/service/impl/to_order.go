package impl

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/shop/order/models"
	"strings"
	"time"
)

func (o *OrderServiceImpl) fillOrderSyncRequestFromDB(c *gin.Context, req *models.OrderSyncRequest) error {
	if req == nil || len(req.Orders) == 0 {
		return errors.New("orders不能为空")
	}

	orders := make([]*models.OrderSyncOrder, 0, len(req.Orders))
	seen := make(map[string]struct{}, len(req.Orders))
	for _, info := range req.Orders {
		if info == nil {
			continue
		}
		tid := strings.TrimSpace(info.Tid)
		info.Tid = tid
		if tid == "" {
			return errors.New("tid不能为空")
		}
		if _, ok := seen[tid]; ok {
			continue
		}
		seen[tid] = struct{}{}

		orderInfo, err := o.orderDao.GetByTid(c, tid)
		if err != nil {
			return err
		}
		if orderInfo == nil {
			return fmt.Errorf("订单不存在: %s", tid)
		}
		orders = append(orders, toOrderSyncOrder(orderInfo))
	}
	if len(orders) == 0 {
		return errors.New("orders不能为空")
	}
	req.Orders = orders
	return nil
}

func toOrderSyncOrder(orderInfo *models.Order) *models.OrderSyncOrder {
	if orderInfo == nil {
		return nil
	}
	return &models.OrderSyncOrder{
		Tid:              orderInfo.Tid,
		Weight:           orderInfo.Weight,
		Size:             orderInfo.Size,
		BuyerNick:        orderInfo.BuyerNick,
		BuyerMessage:     orderInfo.BuyerMessage,
		SellerMemo:       orderInfo.SellerMemo,
		Total:            orderInfo.Total,
		Privilege:        orderInfo.Privilege,
		PostFee:          orderInfo.PostFee,
		ReceiverName:     orderInfo.ReceiverName,
		ReceiverState:    firstNonEmpty(orderInfo.ReceiverProvinceName, orderInfo.ReceiverProvince),
		ReceiverCity:     firstNonEmpty(orderInfo.ReceiverCityName, orderInfo.ReceiverCity),
		ReceiverDistrict: firstNonEmpty(orderInfo.ReceiverDistrictName, orderInfo.ReceiverDistrict),
		ReceiverAddress:  orderInfo.ReceiverAddress,
		ReceiverPhone:    orderInfo.ReceiverPhone,
		ReceiverMobile:   orderInfo.ReceiverMobile,
		Created:          formatOrderSyncTime(orderInfo.CreateTime),
		Status:           orderInfo.Status,
		Type:             orderInfo.Type,
		InvoiceName:      orderInfo.InvoiceName,
		SellerFlag:       orderInfo.SellerFlag,
		PayTime:          formatOrderSyncTime(orderInfo.PayTime),
		LogistBTypeCode:  orderInfo.LogistBTypeCode,
		LogistBillCode:   orderInfo.LogistBillCode,
		BTypeCode:        orderInfo.BTypeCode,
		Details:          toOrderSyncDetails(orderInfo.Details),
		Accounts:         toOrderSyncAccounts(orderInfo.Accounts),
	}
}

func toOrderSyncDetails(details []*models.OrderDetail) []*models.OrderSyncDetail {
	if len(details) == 0 {
		return []*models.OrderSyncDetail{}
	}
	result := make([]*models.OrderSyncDetail, 0, len(details))
	for _, detail := range details {
		if detail == nil {
			continue
		}
		result = append(result, &models.OrderSyncDetail{
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
	return result
}

func toOrderSyncAccounts(accounts []*models.OrderAccount) []*models.OrderSyncAccount {
	if len(accounts) == 0 {
		return []*models.OrderSyncAccount{}
	}
	result := make([]*models.OrderSyncAccount, 0, len(accounts))
	for _, account := range accounts {
		if account == nil {
			continue
		}
		result = append(result, &models.OrderSyncAccount{
			FinanceCode: account.FinanceCode,
			Total:       account.Total,
		})
	}
	return result
}

func formatOrderSyncTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format("2006-01-02 15:04:05")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
