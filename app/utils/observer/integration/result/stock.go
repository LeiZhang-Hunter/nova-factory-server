package result

// StockGetResponse 库存查询响应接口，对应 emall.stock.get 返回。
type StockGetResponse interface {
	GetCode() int64
	GetMessage() string
	GetTotal() int64
	GetStocks() []StockGetResponseData
}

// StockGetResponseData 单条库存数据接口，对应 stocks[{goodsid, ..., whsname}]。
type StockGetResponseData interface {
	GetGoodsid() int64
	GetGoodscode() string
	GetGoodsname() string
	GetSkuid() string
	GetSkucode() string
	GetSkuname() string
	GetQty() float64
	GetEnableNum() float64
	GetEnSaleNum() float64
	GetWhsCode() string
	GetWhsName() string
}
