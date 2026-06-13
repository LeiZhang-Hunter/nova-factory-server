package shopobserver

import (
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/observer"
	"nova-factory-server/app/utils/observer/integration/result"
)

type ShopObserver struct {
	goodsService shopservice.IShopGoodsService
}

func NewShopObserver(goodsService shopservice.IShopGoodsService) *ShopObserver {
	s := &ShopObserver{
		goodsService: goodsService,
	}
	observer.GetNotifier().Register(s)
	return s
}

func (s *ShopObserver) Name() kind.Kind {
	return "shop"
}

func (s *ShopObserver) OnProductChanged(ev event.ProductEvent) (result.SyncProductResponse, error) {
	s.goodsService.SyncEvent(ev)
	return nil, nil
}

func (s *ShopObserver) OnStockChanged(ev event.StockEvent) error {
	return s.goodsService.SyncStock(ev.DB(), ev.GetStocks())
}

func (s *ShopObserver) OnOrderChanged(ev event.OrderEvent) error {
	return nil
}
