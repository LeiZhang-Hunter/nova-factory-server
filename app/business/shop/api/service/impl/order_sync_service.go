//go:build erp
// +build erp

package impl

import (
	"context"
	"errors"
	"fmt"
	erporderdao "nova-factory-server/app/business/erp/sale/saledao"
	erpordermodels "nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"
	orderConstant "nova-factory-server/app/constant/order"
	shopConstant "nova-factory-server/app/constant/shop"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/observer"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	orderSyncLockPrefix = "shop:order:sync:"
	orderSyncLockTTL    = 30 * time.Second
)

type shopOrderSyncService struct {
	cache           cache.Cache
	db              *gorm.DB
	orderDao        erporderdao.IOrderDao
	erpOrderService saleservice.IOrderService
}

func NewShopOrderSyncService(
	cache cache.Cache,
	db *gorm.DB,
	orderDao erporderdao.IOrderDao,
	erpOrderService saleservice.IOrderService,
) *shopOrderSyncService {
	return &shopOrderSyncService{
		cache:           cache,
		db:              db,
		orderDao:        orderDao,
		erpOrderService: erpOrderService,
	}
}

func (s *shopOrderSyncService) SyncCreatedOrder(c *gin.Context, tid string) error {
	tid = strings.TrimSpace(tid)
	if tid == "" {
		return errors.New("订单号不能为空")
	}
	if s == nil || s.orderDao == nil || s.erpOrderService == nil {
		return errors.New("订单同步服务未初始化")
	}

	lockKey := orderSyncLockPrefix + tid
	if s.cache != nil && !s.cache.SetNX(context.Background(), lockKey, "1", orderSyncLockTTL) {
		return errors.New("订单正在同步，请稍后重试")
	}
	if s.cache != nil {
		defer s.cache.Del(context.Background(), lockKey)
	}

	order, err := s.orderDao.GetByTid(c, tid)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("订单不存在")
	}
	if order.SyncStatus == shopConstant.OrderSyncStatusSuccess {
		return nil
	}
	if strings.TrimSpace(order.Status) != orderConstant.ERPStatusNoPay {
		return fmt.Errorf("订单状态不允许同步: %s", order.Status)
	}
	err = observer.GetNotifier().OnOrderChanged(&erpordermodels.OrderSyncRequest{
		Orders: []*erpordermodels.OrderSyncOrder{s.toSyncOrder(order)},
	})
	if err != nil {
		// todo
		return err
	}
	//if err != nil {
	//	_ = s.markSyncFailed(c, order.ID, err.Error())
	//	return err
	//}
	//
	result := findOrderSyncResult(resp, tid)
	if result == nil {
		err = errors.New("管家婆未返回订单同步结果")
		_ = s.markSyncFailed(c, order.ID, err.Error())
		return err
	}
	if err := s.markSyncSuccess(c, order.ID, result.BillCode, result.Message); err != nil {
		return err
	}
	return nil
}

func (s *shopOrderSyncService) markSyncSuccess(c *gin.Context, orderID uint64, billCode string, message string) error {
	now := time.Now()
	if strings.TrimSpace(message) == "" {
		message = "success"
	}
	return s.updateOrderSyncFields(c, orderID, map[string]interface{}{
		"bill_code":    strings.TrimSpace(billCode),
		"sync_status":  shopConstant.OrderSyncStatusSuccess,
		"sync_message": truncateSyncMessage(message),
		"sync_time":    &now,
		"update_time":  &now,
	})
}

func (s *shopOrderSyncService) markSyncFailed(c *gin.Context, orderID uint64, message string) error {
	now := time.Now()
	return s.updateOrderSyncFields(c, orderID, map[string]interface{}{
		"sync_status":  shopConstant.OrderSyncStatusFailed,
		"sync_message": truncateSyncMessage(message),
		"sync_time":    &now,
		"update_time":  &now,
	})
}

func (s *shopOrderSyncService) updateOrderSyncFields(c *gin.Context, orderID uint64, updates map[string]interface{}) error {
	if s == nil || s.db == nil {
		return errors.New("数据库连接不存在")
	}
	return s.db.WithContext(c).Table("erp_order").
		Where("id = ?", orderID).
		Updates(updates).Error
}

func (s *shopOrderSyncService) toSyncOrder(order *erpordermodels.Order) *erpapi.OrderSyncOrder {
	if order == nil {
		return nil
	}
	created := formatOrderTime(order.CreateTime)
	var payTime *string
	if order.PayTime != nil {
		value := order.PayTime.Format("2006-01-02 15:04:05")
		payTime = &value
	}
	return &erpapi.OrderSyncOrder{
		Tid:              order.Tid,
		Weight:           floatPtrIfPositive(order.Weight),
		Size:             floatPtrIfPositive(order.Size),
		BuyerNick:        stringPtrIfNotEmpty(order.BuyerNick),
		BuyerMessage:     stringPtrIfNotEmpty(order.BuyerMessage),
		SellerMemo:       stringPtrIfNotEmpty(order.SellerMemo),
		Total:            floatPtr(order.Total),
		Privilege:        order.Privilege,
		PostFee:          order.PostFee,
		ReceiverName:     order.ReceiverName,
		ReceiverState:    firstNonEmptyString(order.ReceiverProvinceName, order.ReceiverProvince),
		ReceiverCity:     firstNonEmptyString(order.ReceiverCityName, order.ReceiverCity),
		ReceiverDistrict: firstNonEmptyString(order.ReceiverDistrictName, order.ReceiverDistrict),
		ReceiverAddress:  order.ReceiverAddress,
		ReceiverPhone:    stringPtrIfNotEmpty(order.ReceiverPhone),
		ReceiverMobile:   order.ReceiverMobile,
		Created:          created,
		Status:           order.Status,
		Type:             order.Type,
		InvoiceName:      stringPtrIfNotEmpty(order.InvoiceName),
		SellerFlag:       stringPtrIfNotEmpty(order.SellerFlag),
		PayTime:          payTime,
		LogistBTypeCode:  stringPtrIfNotEmpty(order.LogistBTypeCode),
		LogistBillCode:   stringPtrIfNotEmpty(order.LogistBillCode),
		BTypeCode:        stringPtrIfNotEmpty(order.BTypeCode),
		Details:          toSyncDetails(order.Details),
		Accounts:         toSyncAccounts(order.Accounts),
	}
}

func toSyncDetails(details []*erpordermodels.OrderDetail) []*erpapi.OrderSyncDetail {
	rows := make([]*erpapi.OrderSyncDetail, 0, len(details))
	for _, detail := range details {
		if detail == nil {
			continue
		}
		rows = append(rows, &erpapi.OrderSyncDetail{
			OID:            detail.OID,
			Barcode:        stringPtrIfNotEmpty(detail.Barcode),
			EShopGoodsID:   stringPtrIfNotEmpty(detail.EShopGoodsID),
			OuterIID:       stringPtrIfNotEmpty(detail.OuterIID),
			EShopGoodsName: detail.EShopGoodsName,
			EShopSKUId:     stringPtrIfNotEmpty(detail.EShopSkuID),
			EShopSKUName:   stringPtrIfNotEmpty(detail.EShopSkuName),
			NumIID:         int64PtrIfPositive(detail.NumIID),
			SKUId:          int64PtrIfPositive(detail.SkuID),
			Num:            detail.Num,
			Payment:        detail.Payment,
			PicPath:        stringPtrIfNotEmpty(detail.PicPath),
			Weight:         floatPtrIfPositive(detail.Weight),
			Size:           floatPtrIfPositive(detail.Size),
			UnitID:         int64PtrIfPositive(detail.UnitID),
			UnitQty:        floatPtrIfPositive(detail.UnitQty),
		})
	}
	return rows
}

func toSyncAccounts(accounts []*erpordermodels.OrderAccount) []*erpapi.OrderSyncAccount {
	rows := make([]*erpapi.OrderSyncAccount, 0, len(accounts))
	for _, account := range accounts {
		if account == nil {
			continue
		}
		rows = append(rows, &erpapi.OrderSyncAccount{
			FinanceCode: account.FinanceCode,
			Total:       account.Total,
		})
	}
	return rows
}

func findOrderSyncResult(resp *erpapi.OrderSyncResponse, tid string) *erpapi.OrderSyncResult {
	if resp == nil {
		return nil
	}
	tid = strings.TrimSpace(tid)
	for _, result := range resp.Orders {
		if result != nil && strings.TrimSpace(result.Tid) == tid {
			return result
		}
	}
	if len(resp.Orders) == 1 {
		return resp.Orders[0]
	}
	return nil
}

func formatOrderTime(t *time.Time) string {
	if t == nil {
		return time.Now().Format("2006-01-02 15:04:05")
	}
	return t.Format("2006-01-02 15:04:05")
}

func stringPtrIfNotEmpty(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func floatPtr(value float64) *float64 {
	return &value
}

func floatPtrIfPositive(value float64) *float64 {
	if value <= 0 {
		return nil
	}
	return &value
}

func int64PtrIfPositive(value int64) *int64 {
	if value <= 0 {
		return nil
	}
	return &value
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func truncateSyncMessage(message string) string {
	message = strings.TrimSpace(message)
	if len(message) <= 500 {
		return message
	}
	return message[:500]
}
