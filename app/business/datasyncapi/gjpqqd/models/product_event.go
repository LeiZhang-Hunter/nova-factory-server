// 管家婆全渠道商品同步相关的数据模型。
// 定义商品同步请求（GoodsSyncReq）、商品信息（GoodsInfoReq）及 SKU 信息（GoodsSkuSyncReq），
// 同时实现 event.ProductEvent 和 event.ProductData 接口，使其能作为观察者模式的事件载体。
package models

import (
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"

	"gorm.io/gorm"
)

// GoodsSyncReq 商品同步请求，实现 event.ProductEvent 接口
// 作为观察者模式中的事件载体，包含待同步的商品列表、缓存和回调
type GoodsSyncReq struct {
	GoodsInfos []GoodsInfoReq `json:"goodsinfos" form:"goodsinfos"`
	cache      cache.Cache    `json:"-"`
	callback   event.Callback `json:"-"`
	db         *gorm.DB
}

func (g *GoodsSyncReq) GetDB() *gorm.DB {
	return g.db
}

// WidthDB 注入事务 DB，供 Notifier 在事务内分发事件时调用。
func (g *GoodsSyncReq) WidthDB(db *gorm.DB) {
	g.db = db
}

// ToEvent 脱壳为商品变更业务事件，交由观察者处理。
func (g *GoodsSyncReq) ToEvent() event.ProductEvent {
	return g
}

// GetProducts 实现 event.ProductEvent 接口，将内部 GoodsInfoReq 转换为 event.ProductData 列表
func (g *GoodsSyncReq) GetProducts() []event.ProductData {
	goods := make([]event.ProductData, 0)
	//for _, info := range g.GoodsInfos {
	//	//goods = append(goods, &info)
	//}
	return goods
}

// Metadata 实现 event.Base 接口，返回扩展元数据
func (g *GoodsSyncReq) Metadata() map[string]any {
	return map[string]any{}
}

// Ptr 实现 event.Base 接口，返回自身指针
func (g *GoodsSyncReq) Ptr() any {
	return g
}

// Config 实现 event.Event 接口，返回集成配置（当前返回 nil）
func (g *GoodsSyncReq) Config() config.Config {
	return nil
}

// Action 实现 event.Event 接口，商品同步统一视为商品更新事件
func (g *GoodsSyncReq) Action() event.EventType {
	return event.EventProductUpdated
}

// Cache 实现 event.Event 接口，返回缓存实例
func (g *GoodsSyncReq) Cache() cache.Cache {
	return g.cache
}

// WIthCallback 设置事件处理完成后的回调
func (g *GoodsSyncReq) WIthCallback(callback event.Callback) {
	g.callback = callback
	return
}

// GetCallback 实现 event.Event 接口，返回事件回调
func (g *GoodsSyncReq) GetCallback() event.Callback {
	return g.callback
}

// GoodsInfoReq 单个商品信息，实现 event.ProductData 接口
// 包含商品的基本字段（ID、名称、编码、价格、库存等）及 SKU 列表
type GoodsInfoReq struct {
	Goodsid   int64             `json:"goodsid" form:"goodsid"`
	Goodsname string            `json:"goodsname" form:"goodsname"`
	Goodscode string            `json:"goodscode" form:"goodscode"`
	Image     string            `json:"image" form:"image"`
	Desc      string            `json:"desc" form:"desc"`
	Unitname  string            `json:"unitname" form:"unitname"`
	Isservice string            `json:"isservice" form:"isservice"`
	Price     string            `json:"price" form:"price"`
	Quantity  string            `json:"quantity" form:"quantity"`
	Outerid   string            `json:"outerid" form:"outerid"`
	Skus      []GoodsSkuSyncReq `json:"skus" form:"skus"`
}

func (g *GoodsInfoReq) DB() *gorm.DB {
	//TODO implement me
	panic("implement me")
}

// 实现 event.Base 接口
func (g *GoodsInfoReq) Metadata() map[string]any {
	return map[string]any{}
}
func (g *GoodsInfoReq) Ptr() any {
	return g
}

// 实现 event.ProductData 接口的各 getter 方法
func (g *GoodsInfoReq) GetGoodsId() int64 {
	return g.Goodsid
}
func (g *GoodsInfoReq) GetGoodsName() string {
	return g.Goodsname
}
func (g *GoodsInfoReq) GetGoodsCode() string {
	return g.Goodscode
}
func (g *GoodsInfoReq) GetImage() string {
	return g.Image
}
func (g *GoodsInfoReq) GetDesc() string {
	return g.Desc
}
func (g *GoodsInfoReq) GetUnitName() string {
	return g.Unitname
}
func (g *GoodsInfoReq) GetIsService() string {
	return g.Isservice
}
func (g *GoodsInfoReq) GetPrice() string {
	return g.Price
}
func (g *GoodsInfoReq) GetQuantity() string {
	return g.Quantity
}
func (g *GoodsInfoReq) GetOuterId() string {
	return g.Outerid
}

// GetSkus 实现 event.ProductData 接口，返回 SKU 列表
func (g *GoodsInfoReq) GetSkus() []event.ProductSku {
	skus := make([]event.ProductSku, 0)
	return skus
}

// GoodsSkuSyncReq 商品 SKU 同步信息，实现 event.ProductSku 接口
// 描述一个商品的单个规格变体，如不同颜色、尺码的独立 SKU
type GoodsSkuSyncReq struct {
	Goodsid  int64   `json:"goodsid" form:"goodsid"`
	Outerid  string  `json:"outerid" form:"outerid"`
	Price    float64 `json:"price" form:"price"`
	Quantity int64   `json:"quantity" form:"quantity"`
	Skuid    int64   `json:"skuid" form:"skuid"`
	Skuname  string  `json:"skuname" form:"skuname"`
	Skucode  string  `json:"skucode" form:"skucode"`
	Barcode  string  `json:"barcode" form:"barcode"`
	Weight   float64 `json:"weight" form:"weight"`
	Size     int64   `json:"size" form:"size"`
}

// 实现 event.ProductSku 接口的各 getter 方法
func (g *GoodsSkuSyncReq) GetGoodsId() int64 {
	return g.Goodsid
}
func (g *GoodsSkuSyncReq) GetOuterId() string {
	return g.Outerid
}
func (g *GoodsSkuSyncReq) GetPrice() float64 {
	return g.Price
}
func (g *GoodsSkuSyncReq) GetQuantity() int64 {
	return g.Quantity
}
func (g *GoodsSkuSyncReq) GetSkuId() int64 {
	return g.Skuid
}
func (g *GoodsSkuSyncReq) GetSkuName() string {
	return g.Skuname
}
func (g *GoodsSkuSyncReq) GetSkuCode() string {
	return g.Skucode
}
func (g *GoodsSkuSyncReq) GetBarcode() string {
	return g.Barcode
}
func (g *GoodsSkuSyncReq) GetWeight() float64 {
	return g.Weight
}
func (g *GoodsSkuSyncReq) GetSize() int64 {
	return g.Size
}

//type SyncProductResponse struct {
//	IsError  bool   `json:"iserror"`
//	ErrorMsg string `json:"errormsg"`
//	Result   []struct {
//		GoodsId  int    `json:"goodsid"`
//		ErrorMsg string `json:"errormsg"`
//	} `json:"result"`
//}
