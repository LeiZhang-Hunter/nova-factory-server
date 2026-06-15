package event

// ZBtypeGetReqEvent 往来单位查询请求事件，对应 emall.btype.get。
type ZBtypeGetReqEvent struct {
	Event
	Base
	BtypeGetReqDataEvent
}

// BtypeGetReqDataEvent 往来单位查询数据载体。
type BtypeGetReqDataEvent interface {
	GetPage() int64
	GetPageSize() int64
	GetBtypeCode() *string
	GetBtypeName() *string
	GetTel() *string
}
