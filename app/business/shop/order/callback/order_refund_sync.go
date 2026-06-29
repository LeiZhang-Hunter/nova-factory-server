package callback

import (
	"errors"
	"time"

	"nova-factory-server/app/business/shop/order/dao"
	"nova-factory-server/app/business/shop/order/models"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"
	"nova-factory-server/app/utils/store/integration"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AfterSaleSyncCallback 售后同步回调，对标 OrderOrderSyncRequestCallback。
type AfterSaleSyncCallback struct {
	ctx            *gin.Context
	orderRefundDao dao.IOrderRefundDao
	aftersaleID    int64
	event          *models.AftersaleSyncEvent
	isErr          bool
}

func NewAfterSaleSyncCallback(c *gin.Context, orderRefundDao dao.IOrderRefundDao, aftersaleID int64, event *models.AftersaleSyncEvent) *AfterSaleSyncCallback {
	return &AfterSaleSyncCallback{
		ctx:            c,
		orderRefundDao: orderRefundDao,
		aftersaleID:    aftersaleID,
		event:          event,
	}
}

func (s *AfterSaleSyncCallback) OnSuccess(T event.Event, response result.SyncProductResponse) {}

func (s *AfterSaleSyncCallback) OnError(T event.Event, response result.SyncProductResponse, err error) {
	s.isErr = true
}

func (s *AfterSaleSyncCallback) OnFinish(ev event.Event) error {
	if s.isErr {
		return nil
	}
	getService, serviceConfig, err := integration.GetStore().GetService(s.ctx)
	if err != nil {
		zap.L().Error("获取集成商服务失败", zap.Error(err))
		s.updateSyncFailed(err.Error())
		return err
	}
	if getService == nil {
		s.updateSyncFailed("集成商配置不能为空")
		return errors.New("集成商配置不能为空")
	}
	if serviceConfig == nil {
		s.updateSyncFailed("集成商配置不能为空")
		return errors.New("集成商配置不能为空")
	}

	s.event.WithConfig(serviceConfig)
	resp, err := getService.OrderSyncer().SyncAfterSaleOrders(s.ctx, s.event)
	if err != nil {
		zap.L().Error("同步售后单失败", zap.Error(err))
		s.updateSyncFailed(err.Error())
		return err
	}

	now := time.Now()
	billCode := extractAfterSaleBillCode(resp)
	err = s.orderRefundDao.UpdateByID(s.ctx, s.aftersaleID, map[string]any{
		"erp_sync_status":    int32(1),
		"erp_sync_bill_code": billCode,
		"erp_sync_time":      &now,
		"erp_sync_message":   resp.GetMessage(),
	})
	if err != nil {
		zap.L().Error("更新售后单同步状态失败", zap.Error(err))
		return err
	}
	return nil
}

func (s *AfterSaleSyncCallback) updateSyncFailed(message string) {
	now := time.Now()
	s.orderRefundDao.UpdateStatus(s.ctx, s.aftersaleID, 0, map[string]any{
		"sync_status":  int32(2),
		"sync_message": message,
		"sync_time":    &now,
	})
}

func extractAfterSaleBillCode(resp result.AfterSaleOrderSyncResponse) string {
	if resp == nil {
		return ""
	}
	orders := resp.GetOrders()
	if len(orders) == 0 {
		return ""
	}
	for _, o := range orders {
		if o.GetBillCode() != "" {
			return o.GetBillCode()
		}
	}
	return ""
}
