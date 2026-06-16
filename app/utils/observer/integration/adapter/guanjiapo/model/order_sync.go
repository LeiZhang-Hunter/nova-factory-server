package model

import (
	"encoding/json"
	"nova-factory-server/app/utils/observer/integration/result"
)

type OrderSyncResult struct {
	Tid      string `json:"tid"`
	BillCode string `json:"billcode"`
	Message  string `json:"message"`
}

func (o *OrderSyncResult) GetTid() string           { return o.Tid }
func (o *OrderSyncResult) GetBillCode() string      { return o.BillCode }
func (o *OrderSyncResult) GetMessage() string       { return o.Message }
func (o *OrderSyncResult) Ptr() any                 { return o }
func (o *OrderSyncResult) RawStr() (string, error)  { return "", nil }
func (o *OrderSyncResult) MetaData() map[string]any { return make(map[string]any) }

type OrderSyncResponse struct {
	Code    int64             `json:"code"`
	Message string            `json:"message"`
	Orders  []OrderSyncResult `json:"orders"`
}

func (o *OrderSyncResponse) GetCode() int64     { return o.Code }
func (o *OrderSyncResponse) GetMessage() string { return o.Message }
func (o *OrderSyncResponse) GetOrders() []result.OrderSyncResult {
	orders := make([]result.OrderSyncResult, 0, len(o.Orders))
	for _, order := range o.Orders {
		orders = append(orders, &order)
	}
	return orders
}
func (o *OrderSyncResponse) Ptr() any { return o }
func (o *OrderSyncResponse) RawStr() (string, error) {
	marshal, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
func (o *OrderSyncResponse) MetaData() map[string]any { return make(map[string]any) }
