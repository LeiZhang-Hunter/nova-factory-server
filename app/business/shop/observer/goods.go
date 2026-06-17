package observer

import (
	apiService "nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/result"
)

type ShopObserver struct {
	goodsService     shopservice.IShopGoodsService
	orderService     shopservice.IShopOrderService
	orderSendService shopservice.IShopOrderSendService
	apiOrderService  apiService.IApiShopOrderService
}

func NewShopObserver(
	goodsService shopservice.IShopGoodsService,
	orderService shopservice.IShopOrderService,
	orderSendService shopservice.IShopOrderSendService,
	apiOrderService apiService.IApiShopOrderService,
) *ShopObserver {
	return &ShopObserver{
		goodsService:     goodsService,
		orderService:     orderService,
		orderSendService: orderSendService,
		apiOrderService:  apiOrderService,
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
	s.orderSendService.Set(sendEvent)
	return nil
}

// OnOrderStatusChange 订单发货变化
func (o *ShopObserver) OnOrderStatusChange(statusEvent event.ZOrderStatusSyncReqEvent) error {
	err := o.apiOrderService.HandleWechatNotify(statusEvent)
	if err != nil {
		return err
	}
	return nil
}
