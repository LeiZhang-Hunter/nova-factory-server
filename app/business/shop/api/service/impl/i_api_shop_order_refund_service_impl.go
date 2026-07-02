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

	// 加载订单详情
	orderDetails, err := s.orderDetailDao.ListByOrderID(c, shopOrder.ID)
	if err != nil {
		zap.L().Error("获取订单详情失败", zap.Error(err))
		return nil, err
	}
	shopOrder.Details = orderDetails

	// 通过观察者模式触发售后同步（SyncAfterSaleOrder 自行管理事务，ERP 同步回调在事务提交后执行）
	event := models2.NewAftersaleSyncEvent(aftersale, shopOrder)
	cb := callback.NewAfterSaleSyncCallback(c, s.orderRefundDao, aftersale.ID, event)
	event.WithCallback(cb)
	event.WithCache(s.cache)
	event.WithDB(s.db)
	event.WithCtx(c)
	event.WithUserId(userID)

	if err := observer.GetNotifier().OnAfterSaleOrderChanged(event); err != nil {
		zap.L().Error("售后单同步触发失败",
			zap.String("out_refund_no", outRefundNo),
			zap.Error(err),
		)
		return nil, err
	}

	return &apimodels.RefundApplyResp{
		ID:          aftersale.ID,
		OutRefundNo: outRefundNo,
		Status:      orderConstant.AftersaleStatusPendingReview,
		StatusText:  "已提交，等待审核",
	}, nil
}

// SubmitReturnLogistics 提交退货物流信息并再次同步 ERP。
func (s *IApiShopRefundServiceImpl) SubmitReturnLogistics(c *gin.Context, userID int64, req *apimodels.SubmitReturnLogisticsReq) (*apimodels.SubmitReturnLogisticsResp, error) {
	aftersale, err := s.orderRefundDao.GetByOrderId(c, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("查询售后单失败: %v", err)
	}
	if aftersale == nil {
		return nil, errors.New("售后单不存在")
	}
	if aftersale.UserID != userID {
		return nil, errors.New("无权操作此售后单")
	}
	if aftersale.Status != orderConstant.AftersaleStatusApproved {
		return nil, errors.New("当前售后状态不允许提交退货物流")
	}
	if aftersale.PreviousStatus != orderConstant.ERPStatusSended &&
		aftersale.PreviousStatus != orderConstant.ERPStatusPartSend {
		return nil, errors.New("仅退款类型无需退货物流")
	}

	// 更新物流信息
	updates := map[string]any{
		"return_logistics_company": req.LogisticsCompany,
		"return_logistics_code":    req.LogisticsCode,
	}
	if err := s.orderRefundDao.UpdateByID(s.db, aftersale.ID, updates); err != nil {
		return nil, fmt.Errorf("更新退货物流信息失败: %v", err)
	}

	// 再次同步 ERP，携带 logistbillcode
	shopOrder, err := s.orderDao.GetByID(c, uint64(aftersale.OrderID))
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %v", err)
	}
	if shopOrder == nil {
		return nil, errors.New("订单不存在")
	}

	// 加载订单详情
	orderDetails, err := s.orderDetailDao.ListByOrderID(c, shopOrder.ID)
	if err != nil {
		zap.L().Error("获取订单详情失败", zap.Error(err))
		return nil, err
	}
	shopOrder.Details = orderDetails

	event := models2.NewAftersaleSyncEvent(aftersale, shopOrder)
	event.SetLogistBillCode(req.LogisticsCode)
	cb := callback.NewAfterSaleSyncCallback(c, s.orderRefundDao, aftersale.ID, event)
	event.WithCallback(cb)
	event.WithCache(s.cache)
	event.WithDB(s.db)
	event.WithCtx(c)
	event.WithUserId(userID)

	if err := observer.GetNotifier().OnAfterSaleOrderChanged(event); err != nil {
		zap.L().Error("退货物流同步ERP失败",
			zap.String("out_refund_no", aftersale.OutRefundNo),
			zap.Error(err),
		)
		return nil, fmt.Errorf("同步ERP失败: %v", err)
	}

	return &apimodels.SubmitReturnLogisticsResp{
		Status:     aftersale.Status,
		StatusText: orderConstant.GetAftersaleStatusText(aftersale.Status),
	}, nil
}
