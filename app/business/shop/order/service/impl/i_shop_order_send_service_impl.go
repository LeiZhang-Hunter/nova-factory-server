package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/utils/observer/integration/event"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/business/shop/order/service"
	"nova-factory-server/app/constant/commonStatus"
)

// ShopOrderSendServiceImpl 订单发货业务实现。
//
// service 层负责完整业务编排：
// 1. 参数校验；
// 2. 开启事务；
// 3. 写入发货主表；
// 4. 批量写入发货明细；
// 5. 任一步失败时返回 error，触发事务回滚。
type ShopOrderSendServiceImpl struct {
	sendDao   dao.IShopOrderSendDao
	detailDao dao.IShopOrderSendDetailDao
}

// NewShopOrderSendService 创建订单发货服务。
func NewShopOrderSendService(
	sendDao dao.IShopOrderSendDao,
	detailDao dao.IShopOrderSendDetailDao,
) service.IShopOrderSendService {
	return &ShopOrderSendServiceImpl{
		sendDao:   sendDao,
		detailDao: detailDao,
	}
}

// Set 存储订单发货数据（主表 + 明细）。
func (s *ShopOrderSendServiceImpl) Set(sendEvent event.OrderSendEvent) error {
	if sendEvent == nil {
		return errors.New("发货数据不能为空")
	}

	send := models.ToOrderSendByEvent(sendEvent)
	if send == nil {
		return errors.New("发货事件转换失败")
	}

	tid := strings.TrimSpace(send.Tid)
	if tid == "" {
		return errors.New("订单编号 tid 不能为空")
	}
	send.Tid = tid

	execute := func(tx *gorm.DB) error {
		send.State = commonStatus.NORMAL
		if err := s.sendDao.Create(tx, send); err != nil {
			zap.L().Error("订单发货写入失败：创建发货主表失败",
				zap.String("tid", send.Tid), zap.Error(err))
			return err
		}

		if len(send.Details) > 0 {
			if err := s.detailDao.BatchCreate(tx, send.ID, send.Details); err != nil {
				zap.L().Error("订单发货写入失败：创建发货明细失败",
					zap.String("tid", send.Tid),
					zap.Uint64("send_id", send.ID), zap.Error(err))
				return err
			}
		}

		return nil
	}

	if !sendEvent.GetTransaction() {
		return s.sendDao.Transaction(execute)
	}

	return execute(sendEvent.GetDB())

}
