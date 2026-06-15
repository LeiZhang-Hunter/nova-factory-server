package result

// OrderStatusGetResponse 订单状态查询响应接口，对应 emall.orderstatus.get 返回。
type OrderStatusGetResponse interface {
	GetCode() int64
	GetMessage() string
	GetOrderStatus() []OrderStatusGetData
}

// OrderStatusGetData 订单状态数据接口，对应 orderstatus[{eshopbillcode, ..., ordersdetailstatus}]。
type OrderStatusGetData interface {
	GetEshopBillCode() string
	GetStatus() int
	GetLogisticsCode() string
	GetLogisticsName() string
	GetLogistBillCode() string
	GetIsMergeSplit() int
	GetOrdersDetailStatus() []OrderStatusGetDetail
}

// OrderStatusGetDetail 订单明细状态数据接口，对应 ordersdetailstatus[{eshopdetailcode, ..., logistbillcode}]。
type OrderStatusGetDetail interface {
	GetEshopDetailCode() string
	GetStatus() int
	GetLogisticsCode() string
	GetLogisticsName() string
	GetLogistBillCode() string
}
