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
)

type ShopObserver struct {
	goodsService shopservice.IShopGoodsService
	orderService shopservice.IShopOrderService
	skuDao       shopdao.IShopSkuDao
	goodsDao     shopdao.IShopGoodsDao
}

func NewShopObserver(goodsService shopservice.IShopGoodsService,
	orderService shopservice.IShopOrderService,
	skuDao shopdao.IShopSkuDao,
	goodsDao shopdao.IShopGoodsDao,
) *ShopObserver {
	s := &ShopObserver{
		goodsService: goodsService,
		orderService: orderService,
		skuDao:       skuDao,
		goodsDao:     goodsDao,
	}
	observer.GetNotifier().Register(s)
	return s
}

func (s *ShopObserver) Name() kind.Kind {
	return "shop"
}

func (s *ShopObserver) OnProductChanged(ev event.ProductEvent) (result.SyncProductResponse, error) {
	products := ev.GetProducts()
	if len(products) == 0 {
		return nil, nil
	}

	c := &gin.Context{}

	for _, product := range products {
		goodsID := strconv.FormatInt(product.GetGoodsId(), 10)

		goodsUpdates := map[string]any{
			"goods_name":  product.GetGoodsName(),
			"goods_code":  product.GetGoodsCode(),
			"image_url":   product.GetImage(),
			"description": product.GetDesc(),
			"unit":        product.GetUnitName(),
		}

		if err := s.goodsDao.UpsertByGoodsID(c, goodsID, goodsUpdates); err != nil {
			zap.L().Error("ShopObserver: 同步商品失败",
				zap.String("goodsId", goodsID),
				zap.Error(err))
			return nil, fmt.Errorf("同步商品失败 goodsId=%s: %w", goodsID, err)
		}

		var totalQuantity int64
		for _, sku := range product.GetSkus() {
			skuID := strconv.FormatInt(sku.GetSkuId(), 10)

			skuUpdates := map[string]any{
				"goods_id":     goodsID,
				"sku_name":     sku.GetSkuName(),
				"sku_code":     sku.GetSkuCode(),
				"outer_id":     sku.GetOuterId(),
				"barcode":      sku.GetBarcode(),
				"retail_price": sku.GetPrice(),
				"weight":       sku.GetWeight(),
				"quantity":     sku.GetQuantity(),
			}

			if err := s.skuDao.UpsertBySkuID(c, skuID, skuUpdates); err != nil {
				zap.L().Error("ShopObserver: 同步SKU失败",
					zap.String("skuId", skuID),
					zap.Error(err))
				return nil, fmt.Errorf("同步SKU失败 skuId=%s: %w", skuID, err)
			}

			totalQuantity += sku.GetQuantity()
		}

		if err := s.goodsDao.UpdateStockByGoodsID(c, goodsID, totalQuantity); err != nil {
			zap.L().Error("ShopObserver: 更新商品总库存失败",
				zap.String("goodsId", goodsID),
				zap.Int64("quantity", totalQuantity),
				zap.Error(err))
			return nil, fmt.Errorf("更新商品总库存失败 goodsId=%s: %w", goodsID, err)
		}
	}

	return nil, nil
}

func (s *ShopObserver) OnStockChanged(ev event.StockEvent) error {
	stocks := ev.Stocks()
	if len(stocks) == 0 {
		return nil
	}

	c := &gin.Context{}
	goodsIDSet := make(map[string]struct{})

	for _, stock := range stocks {
		skuID := strconv.FormatInt(stock.SkuID(), 10)
		goodsID := strconv.FormatInt(stock.ProductID(), 10)
		afterQty := int64(stock.AfterQty())

		if err := s.skuDao.UpdateStockBySkuID(c, skuID, afterQty); err != nil {
			zap.L().Error("ShopObserver: 更新SKU库存失败",
				zap.String("skuId", skuID),
				zap.Int64("quantity", afterQty),
				zap.Error(err))
			return fmt.Errorf("更新SKU库存失败 skuId=%s: %w", skuID, err)
		}

		goodsIDSet[goodsID] = struct{}{}
	}

	for goodsID := range goodsIDSet {
		totalStock, err := s.skuDao.SumStockByGoodsID(c, goodsID)
		if err != nil {
			zap.L().Error("ShopObserver: 汇总商品库存失败",
				zap.String("goodsId", goodsID),
				zap.Error(err))
			return fmt.Errorf("汇总商品库存失败 goodsId=%s: %w", goodsID, err)
		}

		if err := s.goodsDao.UpdateStockByGoodsID(c, goodsID, totalStock); err != nil {
			zap.L().Error("ShopObserver: 更新商品总库存失败",
				zap.String("goodsId", goodsID),
				zap.Int64("quantity", totalStock),
				zap.Error(err))
			return fmt.Errorf("更新商品总库存失败 goodsId=%s: %w", goodsID, err)
		}
	}

	return nil
}

func (s *ShopObserver) OnOrderChanged(event event.OrderEvent) error {
	return nil
}
