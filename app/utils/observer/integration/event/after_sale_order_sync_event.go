package event

// ZAfterSaleOrderSyncReqEvent 售后订单同步请求事件，对应 emall.afterorder.synchronize。
type ZAfterSaleOrderSyncReqEvent interface {
	Event
	Base
	ZAfterSaleOrderSyncReqDataEvent
}

// ZAfterSaleOrderSyncReqDataEvent 售后订单同步的数据载体。
type ZAfterSaleOrderSyncReqDataEvent interface {
	GetOrders() *[]ZAfterSaleOrderSyncReqData
}

// ZAfterSaleOrderSyncReqData 单条售后订单同步数据。
type ZAfterSaleOrderSyncReqData interface {
	GetRtid() string
	GetTid() string
	GetTotal() float64
	GetPrivilege() float64
	GetPostFee() float64
	GetCreated() string
	GetAftSaleType() string
	GetReasonCode() string
	GetLogistBillCode() string
	GetAftSaleRemark() string
	GetDetails() *[]ZAfterSaleOrderDetail
	GetExDetails() *[]ZAfterSaleOrderExDetail
}

// ZAfterSaleOrderDetail 售后订单退货/退款明细，对应 details 数组。
type ZAfterSaleOrderDetail interface {
	GetOid() string
	GetEshopGoodsName() string
	GetEshopSkuName() string
	GetBackQty() float64
	GetBackTotal() float64
	GetOuterIid() string
}

// ZAfterSaleOrderExDetail 售后订单换货明细，对应 exdetails 数组。
type ZAfterSaleOrderExDetail interface {
	GetOid() string
	GetEshopGoodsName() string
	GetEshopSkuName() string
	GetExchangeQty() float64
	GetBackTotal() float64
	GetOuterIid() string
}
