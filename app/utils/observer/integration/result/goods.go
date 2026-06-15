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
	GetGoodsid() int
	GetEshopgoodsid() string
	GetGoodscode() string
	GetGoodsname() string
	GetUnit() string
	GetRemark() string
	GetSkus() []GoodsGetResponseDataSku
	GetUnits() []GoodsGetResponseDataUnit
}
type GoodsGetResponseDataSku interface {
	GetSkuid() int
	GetSkucode() string
	GetSkuname() string
	GetBarcode() string
	GetLcmccode() string
	GetWeight() int
	GetSize() int
	GetPrice() int
	GetPrice2() int
	GetPrice3() int
	GetPrice4() int
	GetPrice5() int
}

type GoodsGetResponseDataUnit interface {
	GetUnitname() string
	GetBarcode() string
	GetRate() int
	GetPrice() int
	GetPrice2() int
	GetPrice3() int
	GetPrice4() int
	GetPrice5() int
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
