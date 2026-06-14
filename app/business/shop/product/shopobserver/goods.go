package shopobserver

import (
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/result"
)

type ShopObserver struct {
	goodsService shopservice.IShopGoodsService
}

func NewShopObserver(goodsService shopservice.IShopGoodsService) *ShopObserver {
	return &ShopObserver{
		goodsService: goodsService,
	}
}

func (s *ShopObserver) Name() kind.Kind {
	return "shop"
}

// OnProductChanged 商品变更回调，当商品创建、更新或删除时触发
func (s *ShopObserver) OnProductChanged(event event.ProductEvent) (result.SyncProductResponse, error) {
	s.goodsService.SyncEvent(event)
	return nil, nil
}

// OnStockChanged 库存变更回调，当库存数量发生变化时触发
func (s *ShopObserver) OnStockChanged(event event.StockEvent) error {
	return nil
}

// OnOrderChanged 订单变更回调，当订单创建或状态变更（付款、发货等）时触发
func (s *ShopObserver) OnOrderChanged(event event.OrderEvent) error {
	return nil
}
