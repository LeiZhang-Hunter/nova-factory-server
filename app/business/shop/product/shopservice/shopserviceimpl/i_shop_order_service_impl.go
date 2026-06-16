package shopserviceimpl

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/constant/commonStatus"
	objectutil "nova-factory-server/app/utils/json"
	"nova-factory-server/app/utils/observer/integration/event"
	timeutil "nova-factory-server/app/utils/time"
	"strings"
	"time"

	"gorm.io/gorm"
)

// IShopOrderServiceImpl 商城订单同步服务实现。
//
// service 层负责完整业务编排：
// 1. 把 observer 的 event.OrderEvent 转换成商城订单模型；
// 2. 开启事务；
// 3. 按 tid 判断订单新增或更新；
// 4. 判断订单状态是否允许覆盖，避免状态倒退；
// 5. 组合调用订单主表 DAO、明细 DAO、账户 DAO；
// 6. 任一步失败时返回 error，触发事务回滚。
//
// DAO 层只保留单表操作，不再承载同步业务逻辑。
type IShopOrderServiceImpl struct {
	orderDao   shopdao.IShopOrderDao
	detailDao  shopdao.IShopOrderDetailDao
	accountDao shopdao.IShopOrderAccountDao
}

// NewIShopOrderServiceImpl 创建商城订单同步服务。
func NewIShopOrderServiceImpl(
	orderDao shopdao.IShopOrderDao,
	detailDao shopdao.IShopOrderDetailDao,
	accountDao shopdao.IShopOrderAccountDao,
) shopservice.IShopOrderService {
	return &IShopOrderServiceImpl{
		orderDao:   orderDao,
		detailDao:  detailDao,
		accountDao: accountDao,
	}
}

// SyncOrder 同步订单事件到商城订单表。
//
// 该方法是订单同步业务入口。它先把事件转换为 shopmodels.Order，
// 再把订单集合交给 syncOrders 做事务内增量同步。
func (i *IShopOrderServiceImpl) SyncOrder(event event.OrderEvent) error {
	return i.syncOrders(shopmodels.ToOrder(event))
}

// syncOrders 增量同步订单集合。
//
// 同步是“按本次传入数据增量处理”，不会删除本次未传入的其他订单。
// 所有订单共用一个事务：只要任意订单失败，整个批次都会回滚。
// 如果后续业务要求“单个订单失败不影响其它订单”，应把 Transaction 移到循环内部。
func (i *IShopOrderServiceImpl) syncOrders(orders []*shopmodels.Order) error {
	if len(orders) == 0 {
		return nil
	}

	return i.orderDao.Transaction(func(tx *gorm.DB) error {
		for _, item := range orders {
			if item == nil {
				continue
			}
			if err := i.syncOne(tx, item); err != nil {
				zap.L().Error("商城订单同步失败", zap.String("tid", item.Tid), zap.Error(err))
				return err
			}
		}
		return nil
	})
}

// syncOne 在外层事务中同步单个订单。
//
// 处理顺序：
// 1. 校验并标准化 tid；
// 2. 查询 shop_order 是否存在有效记录；
// 3. 准备主表数据，包括 details_json/accounts_json 快照、创建/更新时间、state；
// 4. 不存在则插入主表；
// 5. 已存在则更新主表，删除旧明细和旧账户；
// 6. 插入本次事件携带的新明细和新账户。
//
// 该函数不创建事务，必须使用 syncOrders 传入的 tx。
func (i *IShopOrderServiceImpl) syncOne(tx *gorm.DB, order *shopmodels.Order) error {
	tid := strings.TrimSpace(order.Tid)
	if tid == "" {
		err := errors.New("订单tid不能为空")
		zap.L().Error("商城订单同步失败：订单tid为空", zap.Error(err))
		return err
	}
	order.Tid = tid

	exists, err := i.orderDao.GetByTid(tx, tid)
	if err != nil {
		zap.L().Error("商城订单同步失败：查询已存在订单失败", zap.String("tid", tid), zap.Error(err))
		return err
	}

	now := time.Now()
	if err := prepareOrderForSave(order, &now); err != nil {
		zap.L().Error("商城订单同步失败：准备订单数据失败", zap.String("tid", tid), zap.Error(err))
		return err
	}

	var orderID uint64
	if exists == nil {
		if err := i.orderDao.Create(tx, order); err != nil {
			zap.L().Error("商城订单同步失败：创建订单主表失败", zap.String("tid", tid), zap.Error(err))
			return err
		}
		orderID = order.ID
	} else {
		order.ID = exists.ID
		orderID = exists.ID
		if err := i.orderDao.UpdateByID(tx, exists.ID, buildShopOrderUpdateMap(order, exists)); err != nil {
			zap.L().Error("商城订单同步失败：更新订单主表失败", zap.String("tid", tid), zap.Uint64("order_id", exists.ID), zap.Error(err))
			return err
		}
		if err := i.detailDao.DeleteByOrderID(tx, exists.ID); err != nil {
			zap.L().Error("商城订单同步失败：删除旧订单明细失败", zap.String("tid", tid), zap.Uint64("order_id", exists.ID), zap.Error(err))
			return err
		}
		if err := i.accountDao.DeleteByOrderID(tx, exists.ID); err != nil {
			zap.L().Error("商城订单同步失败：删除旧订单账户失败", zap.String("tid", tid), zap.Uint64("order_id", exists.ID), zap.Error(err))
			return err
		}
	}

	if err := i.detailDao.BatchCreate(tx, orderID, order, &now); err != nil {
		zap.L().Error("商城订单同步失败：创建订单明细失败", zap.String("tid", tid), zap.Uint64("order_id", orderID), zap.Int("details", len(order.Details)), zap.Error(err))
		return err
	}
	if err := i.accountDao.BatchCreate(tx, orderID, order, &now); err != nil {
		zap.L().Error("商城订单同步失败：创建订单账户失败", zap.String("tid", tid), zap.Uint64("order_id", orderID), zap.Int("accounts", len(order.Accounts)), zap.Error(err))
		return err
	}
	return nil
}

// prepareOrderForSave 准备订单主表写入前的数据。
//
// Details 和 Accounts 会序列化成主表 JSON 快照字段，便于只查主表时保留同步载荷。
// 同时补齐 create_time/update_time/state。
func prepareOrderForSave(order *shopmodels.Order, now *time.Time) error {
	detailsJSON, err := objectutil.MarshalJSON(order.Details)
	if err != nil {
		zap.L().Error("商城订单同步失败：订单明细JSON序列化失败", zap.String("tid", order.Tid), zap.Int("details", len(order.Details)), zap.Error(err))
		return fmt.Errorf("订单明细JSON序列化失败: %w", err)
	}
	accountsJSON, err := objectutil.MarshalJSON(order.Accounts)
	if err != nil {
		zap.L().Error("商城订单同步失败：订单账户JSON序列化失败", zap.String("tid", order.Tid), zap.Int("accounts", len(order.Accounts)), zap.Error(err))
		return fmt.Errorf("订单账户JSON序列化失败: %w", err)
	}

	order.DetailsJSON = detailsJSON
	order.AccountsJSON = accountsJSON
	order.CreateTime = timeutil.FirstTime(order.CreateTime, now)
	order.UpdateTime = timeutil.FirstTime(order.UpdateTime, now)
	order.State = commonStatus.NORMAL
	return nil
}

// SyncOrderStatus 同步订单状态
func (i *IShopOrderServiceImpl) SyncOrderStatus(event event.OrderStratusEvent) error {
	if event == nil || len(event.GetOrders()) == 0 {
		return nil
	}
	if event.GetDB() == nil {
		return nil
	}

	now := time.Now()
	for _, item := range event.GetOrders() {
		if item == nil {
			continue
		}

		tid := strings.TrimSpace(item.GetTid())
		status := strings.TrimSpace(item.GetStatus())
		if tid == "" {
			return errors.New("订单tid不能为空")
		}
		if status == "" {
			continue
		}

		order, err := i.orderDao.GetByTid(event.GetDB(), tid)
		if err != nil {
			zap.L().Error("商城订单状态同步失败：查询订单失败", zap.String("tid", tid), zap.Error(err))
			return err
		}
		if order == nil {
			zap.L().Debug("商城订单状态同步跳过不存在订单", zap.String("tid", tid), zap.String("status", status))
			continue
		}
		if !shouldUpdateOrderStatus(order.Status, status) {
			zap.L().Debug("商城订单状态同步跳过状态更新",
				zap.String("tid", tid),
				zap.String("current_status", order.Status),
				zap.String("incoming_status", status),
			)
			continue
		}

		if err := i.orderDao.UpdateByID(event.GetDB(), order.ID, map[string]any{
			"status":      status,
			"update_time": &now,
		}); err != nil {
			zap.L().Error("商城订单状态同步失败：更新订单状态失败", zap.String("tid", tid), zap.Uint64("order_id", order.ID), zap.Error(err))
			return err
		}
	}
	return nil
}
