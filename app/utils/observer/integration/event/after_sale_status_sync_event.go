package event

// ZAfterSaleStatusSyncReqEvent ERP售后状态回调事件，对应 selfmall.afterorder.status.sync。
type ZAfterSaleStatusSyncReqEvent interface {
	Event
	Base
	GetRtid() string
	GetTid() string
	GetStatus() string
}
