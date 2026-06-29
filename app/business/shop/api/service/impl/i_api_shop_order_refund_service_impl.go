package impl

import (
	"context"
	"errors"
	"fmt"
	"nova-factory-server/app/business/shop/order/callback"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/observer"
	"nova-factory-server/app/utils/snowflake"
	"time"

	"nova-factory-server/app/business/shop/api/dao"
	apimodels "nova-factory-server/app/business/shop/api/models"
	apiservice "nova-factory-server/app/business/shop/api/service"
	orderDao "nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/business/shop/order/models"
	models2 "nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/business/shop/order/provider"
	orderConstant "nova-factory-server/app/constant/order"
	"nova-factory-server/app/utils/order"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	aftersaleCreateLockPrefix = "shop:app:aftersale:create:"
	aftersaleCreateLockTTL    = 15 * time.Second
)

// IApiShopRefundServiceImpl 小程序售后业务实现。
type IApiShopRefundServiceImpl struct {
	orderDetailDao orderDao.IOrderDetailDao
	orderDao       orderDao.IOrderDao
	configDao      dao.IApiShopSysConfigDao
	orderRefundDao orderDao.IOrderRefundDao
	cache          cache.Cache
	db             *gorm.DB
}

// NewIApiShopRefundService 创建小程序售后业务服务。
func NewIApiShopRefundService(
	orderDao orderDao.IOrderDao,
	orderDetailDao orderDao.IOrderDetailDao,
	configDao dao.IApiShopSysConfigDao,
	cache cache.Cache,
	orderRefundDao orderDao.IOrderRefundDao,
	db *gorm.DB,
) apiservice.IApiShopOrderRefundService {
	return &IApiShopRefundServiceImpl{
		orderDao:       orderDao,
		configDao:      configDao,
		cache:          cache,
		orderRefundDao: orderRefundDao,
		orderDetailDao: orderDetailDao,
		db:             db,
	}
}

func (s *IApiShopRefundServiceImpl) Apply(c *gin.Context, userID int64, req *apimodels.RefundApplyReq) (*apimodels.RefundApplyResp, error) {
	// 查询订单
	shopOrder, err := s.orderDao.GetByID(c, uint64(req.OrderID))
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %v", err)
	}
	if shopOrder == nil {
		return nil, errors.New("订单不存在")
	}
	if shopOrder.UserId != uint64(userID) {
		return nil, errors.New("无权操作此订单")
	}
	// 校验订单状态是否可售后
	if shopOrder.Status == orderConstant.ERPStatusNoPay ||
		shopOrder.Status == orderConstant.ERPStatusTradeClosed ||
		shopOrder.Status == orderConstant.ERPStatusAftersale {
		return nil, errors.New("当前订单状态不可申请售后")
	}
	// Redis 分布式锁，防止重复提交
	lockKey := aftersaleCreateLockPrefix + fmt.Sprintf("%d", req.OrderID)
	if s.cache != nil && !s.cache.SetNX(context.Background(), lockKey, "1", aftersaleCreateLockTTL) {
		return nil, errors.New("请勿重复提交，售后申请处理中")
	}
	defer func() {
		if s.cache != nil {
			s.cache.Del(context.Background(), lockKey)
		}
	}()

	// 检查是否存在进行中的售后单
	existing, err := s.orderRefundDao.GetByOrderId(c, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("查询售后申请失败: %v", err)
	}
	if existing != nil {
		return nil, errors.New("该订单已存在售后申请")
	}
	outRefundNo := order.GenerateOutRefundNo(shopOrder.Tid)
	now := time.Now()
	aftersale := &models.OrderRefund{
		ID:             snowflake.GenID(),
		OrderID:        int64(shopOrder.ID),
		Tid:            shopOrder.Tid,
		UserID:         userID,
		PayChannel:     shopOrder.PayChannel,
		RefundChannel:  shopOrder.PayChannel,
		OutRefundNo:    outRefundNo,
		RefundAmount:   shopOrder.Total,
		TotalAmount:    shopOrder.Total,
		Reason:         req.Reason,
		PreviousStatus: shopOrder.Status,
		Status:         orderConstant.AftersaleStatusPendingReview,
	}
	aftersale.SetCreateBy(userID)

	err = s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := s.orderRefundDao.CreateWithTx(tx, aftersale); err != nil {
			return fmt.Errorf("创建售后申请失败: %v", err)
		}
		return s.orderDao.UpdateByID(tx, shopOrder.ID, map[string]any{
			"status":      orderConstant.ERPStatusAftersale,
			"update_time": &now,
		})
	})
	if err != nil {
		zap.L().Error("更新订单状态失败", zap.Error(err))
		return nil, err
	}

	// 判断是否自动退款
	autoRefund, err := s.configDao.GetIsAutoRefundEnabled(c)
	if err != nil {
		zap.L().Error("获取是否自动退款失败", zap.Error(err))
		return nil, err
	}
	isUnshipped := shopOrder.Status == orderConstant.ERPStatusPayed
	if autoRefund && isUnshipped {
		// 直接发起退款
		refundProvider, err := provider.GetPaymentMethod(shopOrder.PayChannel)
		if err != nil {
			zap.L().Error("获取退款provider失败", zap.Error(err))
			return &apimodels.RefundApplyResp{
				ID:          aftersale.ID,
				OutRefundNo: outRefundNo,
				Status:      orderConstant.AftersaleStatusPendingReview,
				StatusText:  "退款provider不可用，等待人工处理",
			}, nil
		}

		err = s.orderRefundDao.UpdateStatus(c, aftersale.ID, orderConstant.AftersaleStatusRefunding, nil)
		if err != nil {
			zap.L().Error("更新售后申请状态失败", zap.Error(err))
			return nil, err
		}

		provReq := &models.RefundRequest{
			Tid:           shopOrder.Tid,
			TransactionID: shopOrder.TransactionID,
			OutRefundNo:   outRefundNo,
			Reason:        req.Reason,
			RefundAmount:  shopOrder.Total,
			TotalAmount:   shopOrder.Total,
			MchID:         shopOrder.MchID,
			AppID:         shopOrder.AppID,
		}

		refundResult, err := refundProvider.Refund(c, provReq)
		if err != nil {
			zap.L().Error("退款发起失败", zap.Error(err))
			err = s.orderRefundDao.UpdateStatus(c, aftersale.ID, orderConstant.AftersaleStatusRefundFailed, nil)
			if err != nil {
				zap.L().Error("更新售后申请状态失败", zap.Error(err))
				return nil, err
			}
			return nil, errors.New("退款失败")
		}

		updates := map[string]any{
			"third_refund_id": refundResult.ThirdRefundID,
		}
		if refundResult.ThirdTransactionID != "" {
			updates["third_transaction_id"] = refundResult.ThirdTransactionID
		}
		err = s.orderRefundDao.UpdateStatus(c, aftersale.ID, orderConstant.AftersaleStatusRefundSuccess, updates)
		if err != nil {
			zap.L().Error("更新售后申请状态失败", zap.Error(err))
			return nil, err
		}
		// 获取订单详情
		orderDetails, err := s.orderDetailDao.ListByOrderID(c, shopOrder.ID)
		if err != nil {
			zap.L().Error("获取订单详情失败", zap.Error(err))
			return nil, err
		}
		shopOrder.Details = orderDetails
		// 触发管家婆售后同步
		event := models2.NewAftersaleSyncEvent(aftersale, shopOrder)
		cb := callback.NewAfterSaleSyncCallback(c, s.orderRefundDao, aftersale.ID, event)
		event.WithCallback(cb)
		event.WithCache(s.cache)
		event.WithDB(s.db)

		if err := observer.GetNotifier().OnAfterSaleOrderChanged(event); err != nil {
			zap.L().Error("售后单同步触发失败",
				zap.String("out_refund_no", aftersale.OutRefundNo),
				zap.Error(err),
			)
		}
		return &apimodels.RefundApplyResp{
			ID:          aftersale.ID,
			OutRefundNo: outRefundNo,
			Status:      orderConstant.AftersaleStatusRefundSuccess,
			StatusText:  "退款成功",
		}, nil
	}

	return &apimodels.RefundApplyResp{
		ID:          aftersale.ID,
		OutRefundNo: outRefundNo,
		Status:      orderConstant.AftersaleStatusPendingReview,
		StatusText:  "已提交，等待审核",
	}, nil
}
