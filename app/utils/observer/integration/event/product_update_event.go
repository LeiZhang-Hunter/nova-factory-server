package event

type ZProductUpdateReqEvent interface {
	Event
	Base
	ZProductUpdateReqDataEvent
}

type ZProductUpdateReqDataEvent interface {
	GetItems() *[]ZProductUpdateReqData
}
type ZProductUpdateReqData interface {
	GetGoodsID() string
	GetRemark() string
}
