package shopobserver

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/observer"
	"nova-factory-server/app/utils/observer/integration/result"

	"gorm.io/gorm"
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

func (s *ShopObserver) OnProductChanged(tx *gorm.DB, ev event.ProductEvent) (result.SyncProductResponse, error) {
	s.goodsService.SyncEvent(ev)
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
