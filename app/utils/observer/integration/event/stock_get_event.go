package event

// ZStockGetReqEvent 库存查询请求事件，对应 emall.stock.get。
type ZStockGetReqEvent struct {
	Event
	Base
	StockGetReqDataEvent
}

// StockGetReqDataEvent 库存查询数据载体。
type StockGetReqDataEvent interface {
	GetPage() int64
	GetPageSize() int64
	GetSkuCode() *string
	GetGoodsCode() *[]string
	GetWhsCode() *string
	GetIsContainWhs() *bool
}
