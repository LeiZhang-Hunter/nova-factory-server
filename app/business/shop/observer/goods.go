package observer

import (
	"nova-factory-server/app/business/shop/api/models"
	apiService "nova-factory-server/app/business/shop/api/service"
	orderservice "nova-factory-server/app/business/shop/order/service"
	"nova-factory-server/app/business/shop/product/shopservice"
	orderConstant "nova-factory-server/app/constant/order"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/result"

	"go.uber.org/zap"
)

type ShopObserver struct {
	goodsService    shopservice.IShopGoodsService
	orderService    orderservice.IOrderService
	apiOrderService apiService.IApiShopOrderService
}

func NewShopObserver(
	goodsService shopservice.IShopGoodsService,
	orderService orderservice.IOrderService,
	apiOrderService apiService.IApiShopOrderService,
) *ShopObserver {
	return &ShopObserver{
		goodsService:    goodsService,
		orderService:    orderService,
		apiOrderService: apiOrderService,
	}
}

func (s *ShopObserver) Name() kind.Kind {
	return "shop"
}

// OnProductChanged 商品变更回调，当商品创建、更新或删除时触发
func (s *ShopObserver) OnProductChanged(event event.ProductEvent) (result.SyncProductResponse, error) {
	return s.goodsService.SyncEvent(event)
}

// OnStockChanged 库存变更回调，当库存数量发生变化时触发
func (s *ShopObserver) OnStockChanged(event event.StockEvent) error {
	return s.goodsService.SyncStock(event)
}

// OnOrderChanged 订单变更回调，当订单创建或状态变更（付款、发货等）时触发
func (s *ShopObserver) OnOrderChanged(event event.OrderEvent) error {
	if s.orderService == nil {
		return nil
	}
	return s.orderService.SyncOrder(event)
}

// OnOrderSendChange 订单发货变化
func (s *ShopObserver) OnOrderSendChange(sendEvent event.OrderSendEvent) error {
	//s.orderService.Set(sendEvent)
	return nil
}

// OnOrderStatusChange 订单发货变化
func (o *ShopObserver) OnOrderStatusChange(statusEvent event.ZOrderStatusSyncReqEvent) error {
	uid := statusEvent.GetUserId()
	//BatchUpdateStatus(c *gin.Context, userID int64, req *apimodels.BatchOrderStatusReq) error
	var statusList []models.OrderStatus = make([]models.OrderStatus, 0)
	for _, v := range statusEvent.GetOrders() {
		statusList = append(statusList, models.OrderStatus{
			ID:     v.GetDBID(),
			Status: orderConstant.ErpStatusToShopStatus(v.GetStatus()),
		})
	}
	err := o.apiOrderService.BatchUpdateStatus(statusEvent.GetCtx(), uid, statusList)
	if err != nil {
		zap.L().Error("batch update status error", zap.Error(err))
		return err
	}
	return nil
}
