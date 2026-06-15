package event

type ZProductGetReqEvent struct {
	Event
	Base
	ProductGetReqDataEvent
}

type ProductGetReqDataEvent interface {
	GetPage() int64
	GetPageSize() int64
	GetReturnType() *int64
	GetGoodsCode() *[]string
	GetGoodsName() *[]string
}
