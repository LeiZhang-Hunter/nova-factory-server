package impl

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"nova-factory-server/app/business/shop/order/service"
	"nova-factory-server/app/constant/commonStatus"
	orderConstant "nova-factory-server/app/constant/order"
	"nova-factory-server/app/datasource/cache"
)

// OrderTimeoutServiceImpl 订单超时自动取消实现。
//
// 主路径：StartConsumer 启动的常驻 goroutine，每秒轮询 Redis 延迟队列。
// 兜底路径：HTTP 接口调用 ProcessExpiredOrders。
// 两条路径共用 ProcessExpiredOrders，保证取消逻辑一致。
type OrderTimeoutServiceImpl struct {
	cache cache.Cache
	db    *gorm.DB
	table string
}

// NewIOrderTimeoutServiceImpl 创建订单超时自动取消服务。
func NewIOrderTimeoutServiceImpl(cache cache.Cache, db *gorm.DB) service.IOrderTimeoutService {
	return &OrderTimeoutServiceImpl{
		cache: cache,
		db:    db,
		table: "shop_order",
	}
}

// StartConsumer 启动 Consumer goroutine，每秒轮询一次延迟队列。
func (s *OrderTimeoutServiceImpl) StartConsumer(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	zap.L().Info("订单超时取消Consumer已启动")
	for {
		select {
		case <-ctx.Done():
			zap.L().Info("订单超时取消Consumer已停止")
			return
		case <-ticker.C:
			if _, err := s.ProcessExpiredOrders(ctx); err != nil {
				zap.L().Error("订单超时取消轮询失败", zap.Error(err))
			}
		}
	}
}

// ProcessExpiredOrders 取出 score ≤ now 的到期订单并逐个取消，返回成功取消数量。
func (s *OrderTimeoutServiceImpl) ProcessExpiredOrders(ctx context.Context) (int, error) {
	now := time.Now().Unix()
	ids, err := s.cache.ZRangeByScore(ctx, orderConstant.DelayCancelKey, &redis.ZRangeBy{
		Min: "0",
		Max: strconv.FormatInt(now, 10),
	}).Result()
	if err != nil {
		return 0, err
	}

	cancelled := 0
	for _, idStr := range ids {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			s.removeFromQueue(ctx, idStr)
			continue
		}
		id, parseErr := strconv.ParseUint(idStr, 10, 64)
		if parseErr != nil {
			zap.L().Error("订单超时取消跳过非法订单ID", zap.String("member", idStr), zap.Error(parseErr))
			s.removeFromQueue(ctx, idStr)
			continue
		}
		done, cancelErr := s.cancelOne(ctx, id)
		if cancelErr != nil {
			// 处理失败不移除队列成员，等待下次轮询或兜底接口重试。
			zap.L().Error("订单超时取消失败", zap.Uint64("order_id", id), zap.Error(cancelErr))
			continue
		}
		if done {
			cancelled++
		}
		s.removeFromQueue(ctx, idStr)
	}
	return cancelled, nil
}

// cancelOne 取消单个订单。返回是否实际执行了取消。
//
// 幂等：仅当订单存在且为待支付（ERP NoPay）时取消；已支付/已取消/不存在均视为无需取消。
// 取消不绑定 buyer_nick，直接按 id 更新（Decision 5：自动取消场景没有具体用户）。
func (s *OrderTimeoutServiceImpl) cancelOne(ctx context.Context, id uint64) (bool, error) {
	var statusStr string
	err := s.db.WithContext(ctx).
		Table(s.table).
		Select("status").
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Limit(1).
		Scan(&statusStr).Error
	if err != nil {
		return false, err
	}
	if strings.TrimSpace(statusStr) == "" {
		// 订单不存在，无需取消，交由调用方移除队列。
		return false, nil
	}
	if strings.TrimSpace(statusStr) != orderConstant.ERPStatusNoPay {
		// 已支付或已取消等，跳过。
		return false, nil
	}

	result := s.db.WithContext(ctx).
		Table(s.table).
		Where("id = ?", id).
		Where("status = ?", orderConstant.ERPStatusNoPay).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]any{
			"status":      orderConstant.ERPStatusTradeClosed,
			"seller_memo": "超时未支付自动取消",
			"update_time": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// removeFromQueue 从延迟队列移除订单成员，忽略移除错误（仅记录日志）。
func (s *OrderTimeoutServiceImpl) removeFromQueue(ctx context.Context, member string) {
	if err := s.cache.ZRem(ctx, orderConstant.DelayCancelKey, member).Err(); err != nil {
		zap.L().Error("订单超时取消移除队列失败", zap.String("member", member), zap.Error(err))
	}
}
