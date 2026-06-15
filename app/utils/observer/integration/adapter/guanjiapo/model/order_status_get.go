package model

import "nova-factory-server/app/utils/observer/integration/result"

// OrderStatusGetDetail 订单明细状态，实现 result.OrderStatusGetDetail。
type OrderStatusGetDetail struct {
	EshopDetailCode string `json:"eshopdetailcode"`
	Status          int    `json:"status"`
	LogisticsCode   string `json:"logisticscode"`
	LogisticsName   string `json:"logisticsname"`
	LogistBillCode  string `json:"logistbillcode"`
}

func (d *OrderStatusGetDetail) GetEshopDetailCode() string { return d.EshopDetailCode }
func (d *OrderStatusGetDetail) GetStatus() int             { return d.Status }
func (d *OrderStatusGetDetail) GetLogisticsCode() string   { return d.LogisticsCode }
func (d *OrderStatusGetDetail) GetLogisticsName() string   { return d.LogisticsName }
func (d *OrderStatusGetDetail) GetLogistBillCode() string  { return d.LogistBillCode }

// OrderStatusGetData 订单状态数据，实现 result.OrderStatusGetData。
type OrderStatusGetData struct {
	EshopBillCode      string                  `json:"eshopbillcode"`
	Status             int                     `json:"status"`
	LogisticsCode      string                  `json:"logisticscode"`
	LogisticsName      string                  `json:"logisticsname"`
	LogistBillCode     string                  `json:"logistbillcode"`
	IsMergeSplit       int                     `json:"ismergesplit"`
	OrdersDetailStatus []*OrderStatusGetDetail `json:"ordersdetailstatus"`
}

func (d *OrderStatusGetData) GetEshopBillCode() string  { return d.EshopBillCode }
func (d *OrderStatusGetData) GetStatus() int            { return d.Status }
func (d *OrderStatusGetData) GetLogisticsCode() string  { return d.LogisticsCode }
func (d *OrderStatusGetData) GetLogisticsName() string  { return d.LogisticsName }
func (d *OrderStatusGetData) GetLogistBillCode() string { return d.LogistBillCode }
func (d *OrderStatusGetData) GetIsMergeSplit() int      { return d.IsMergeSplit }
func (d *OrderStatusGetData) GetOrdersDetailStatus() []result.OrderStatusGetDetail {
	out := make([]result.OrderStatusGetDetail, len(d.OrdersDetailStatus))
	for i, v := range d.OrdersDetailStatus {
		out[i] = v
	}
	return out
}

// OrderStatusGetResponse 订单状态查询响应，实现 result.OrderStatusGetResponse。
type OrderStatusGetResponse struct {
	Code        int64                 `json:"code"`
	Message     string                `json:"message"`
	OrderStatus []*OrderStatusGetData `json:"orderstatus"`
}

func (r *OrderStatusGetResponse) GetCode() int64     { return r.Code }
func (r *OrderStatusGetResponse) GetMessage() string { return r.Message }
func (r *OrderStatusGetResponse) GetOrderStatus() []result.OrderStatusGetData {
	out := make([]result.OrderStatusGetData, len(r.OrderStatus))
	for i, v := range r.OrderStatus {
		out[i] = v
	}
	return out
}
