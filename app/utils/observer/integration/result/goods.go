package result

type SyncProductResponse interface {
	GetCode() int
	GetMessage() string
	GetMetadata() map[string]any
}

type GoodsGetResponse interface {
	GetCode() int
	GetMessage() string
	GetTotal() int
	GetGoods() []GoodsGetResponseData
}

type GoodsGetResponseData interface {
	GetGoodsid() string
	GetEshopgoodsid() string
	GetGoodscode() string
	GetGoodsname() string
	GetUnit() string
	GetRemark() string
	GetSkus() []GoodsGetResponseDataSku
	GetUnits() []GoodsGetResponseDataUnit
}
type GoodsGetResponseDataSku interface {
	GetSkuid() int64
	GetSkucode() string
	GetSkuname() string
	GetBarcode() string
	GetLcmccode() string
	GetWeight() float64
	GetSize() float64
	GetPrice() float64
	GetPrice2() float64
	GetPrice3() float64
	GetPrice4() float64
	GetPrice5() float64
}

type GoodsGetResponseDataUnit interface {
	GetUnitname() string
	GetBarcode() string
	GetRate() float64
	GetPrice() float64
	GetPrice2() float64
	GetPrice3() float64
	GetPrice4() float64
	GetPrice5() float64
}

// 商品备注更新
type ProductRemarkUpdateResponse interface {
	GetCode() int64
	GetMessage() string
	GetItems() []ProductRemarkUpdateResult
}
type ProductRemarkUpdateResult interface {
	GetGoodsID() int64
	GetMessage() string
}
