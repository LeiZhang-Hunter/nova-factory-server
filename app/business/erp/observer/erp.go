package observer

import (
	"nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/business/erp/sale/saleservice"
	"nova-factory-server/app/business/erp/stock/stockservice"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/kind"
	"nova-factory-server/app/utils/observer/integration/result"
)

// ERPObserver ERP 模块观察者。
//
// 该观察者实现 integration/observer.Observer 接口，用于接收全局 Notifier 分发的业务事件。
// 当前 ERP 侧已接入订单事件：收到订单变更后，会把 event.OrderEvent 转换为 ERP 销售订单保存参数，
// 并调用 saleservice.IOrderService.Set 完成 ERP 订单新增或更新。
//
// 商品和库存事件目前 ERP 模块没有统一的事件同步 service，因此这里保留完整回调入口并记录日志，
// 后续如果 ERP 商品档案或库存单据需要响应事件，可在对应回调中注入并调用相关 service。
type ERPObserver struct {
	orderService   saleservice.IOrderService
	productService masterservice.IProductService
	stockService   stockservice.IStockService
}

// NewERPObserver 创建 ERP 观察者。
func NewERPObserver(orderService saleservice.IOrderService, productService masterservice.IProductService, stockService stockservice.IStockService) *ERPObserver {
	s := &ERPObserver{
		orderService:   orderService,
		productService: productService,
		stockService:   stockService,
	}
	return s
}

// Name 返回观察者名称。
//
// 该名称用于 Notifier 日志和调试时区分不同观察者。
func (o *ERPObserver) Name() kind.Kind {
	return kind.Kind("erp")
}

// OnProductChanged 处理商品变更事件。
//
// 当前 ERP 商品事件同步尚未接入，先安全跳过。
func (o *ERPObserver) OnProductChanged(event event.ProductEvent) (result.SyncProductResponse, error) {
	//products := event.GetProducts()
	//if len(products) == 0 {
	//	return nil, nil
	//}
	//
	//c := &gin.Context{}
	//
	//for _, product := range products {
	//	skus := product.GetSkus()
	//	if len(skus) == 0 {
	//		continue
	//	}
	//
	//	for _, sku := range skus {
	//		skuID := sku.GetSkuId()
	//
	//		productUpdates := map[string]any{
	//			"name":       product.GetGoodsName(),
	//			"bar_code":   sku.GetBarcode(),
	//			"weight":     sku.GetWeight(),
	//			"sale_price": sku.GetPrice(),
	//			"standard":   sku.GetSkuName(),
	//		}
	//		if code := product.GetGoodsCode(); code != "" {
	//			productUpdates["product_code"] = code
	//		}
	//
	//		if err := o.productDao.UpsertByID(c, skuID, productUpdates); err != nil {
	//			zap.L().Error("ErpObserver: 同步ERP产品失败",
	//				zap.Int64("skuId", skuID),
	//				zap.Error(err))
	//			return nil, fmt.Errorf("同步ERP产品失败 skuId=%d: %w", skuID, err)
	//		}
	//
	//		stockUpdates := map[string]any{
	//			"product_id": skuID,
	//			"count":      float64(sku.GetQuantity()),
	//		}
	//
	//		if err := s.stockDao.UpsertByID(c, skuID, stockUpdates); err != nil {
	//			zap.L().Error("ErpObserver: 同步ERP库存失败",
	//				zap.Int64("skuId", skuID),
	//				zap.Error(err))
	//			return nil, fmt.Errorf("同步ERP库存失败 skuId=%d: %w", skuID, err)
	//		}
	//	}
	//}

	return nil, nil
}

// OnStockChanged 处理库存变更事件。
func (o *ERPObserver) OnStockChanged(ev event.StockEvent) error {
	return o.stockService.SyncStock(ev)
}

// OnOrderChanged 处理订单变更事件。
//
// Observer 只负责把事件分发给 ERP 订单 service，不在这里做订单转换和保存。
// 具体同步逻辑由 saleservice.IOrderService.Sync(event event.OrderEvent) 承担。
func (o *ERPObserver) OnOrderChanged(event event.OrderEvent) error {
	//if event == nil {
	//	return nil
	//}
	//orders := event.GetOrders()
	//if len(orders) == 0 {
	//	return nil
	//}
	//if o.orderService == nil {
	//	zap.L().Warn("ERP观察者订单服务为空，跳过订单事件", zap.Int("orders", len(orders)))
	//	return nil
	//}
	//
	//o.orderService.Sync(event)
	return nil
}

// OnOrderSendChange 订单发货变化
func (o *ERPObserver) OnOrderSendChange(sendEvent event.OrderSendEvent) error {
	return nil
}
