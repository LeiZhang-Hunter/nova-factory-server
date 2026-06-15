package model

import (
	"encoding/json"
	"nova-factory-server/app/utils/observer/integration/result"
)

// AfterSaleOrderSyncResult 单笔售后订单同步结果，实现 result.AfterSaleOrderSyncResult。
type AfterSaleOrderSyncResult struct {
	IsError  bool   `json:"iserror"`
	Tid      string `json:"tid"`
	BillCode string `json:"billcode"`
	Message  string `json:"message"`
}

func (r *AfterSaleOrderSyncResult) GetIsError() bool    { return r.IsError }
func (r *AfterSaleOrderSyncResult) GetTid() string      { return r.Tid }
func (r *AfterSaleOrderSyncResult) GetBillCode() string { return r.BillCode }
func (r *AfterSaleOrderSyncResult) GetMessage() string  { return r.Message }
func (r *AfterSaleOrderSyncResult) Ptr() any            { return r }
func (r *AfterSaleOrderSyncResult) RawStr() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func (r *AfterSaleOrderSyncResult) MetaData() map[string]any { return make(map[string]any) }

// AfterSaleOrderSyncResponse 售后订单同步响应，实现 result.AfterSaleOrderSyncResponse。
type AfterSaleOrderSyncResponse struct {
	Code    int64                       `json:"code"`
	Message string                      `json:"message"`
	Orders  []*AfterSaleOrderSyncResult `json:"orders"`
}

func (r *AfterSaleOrderSyncResponse) GetCode() int64     { return r.Code }
func (r *AfterSaleOrderSyncResponse) GetMessage() string { return r.Message }
func (r *AfterSaleOrderSyncResponse) GetOrders() []result.AfterSaleOrderSyncResult {
	out := make([]result.AfterSaleOrderSyncResult, len(r.Orders))
	for i, v := range r.Orders {
		out[i] = v
	}
	return out
}
func (r *AfterSaleOrderSyncResponse) Ptr() any { return r }
func (r *AfterSaleOrderSyncResponse) RawStr() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func (r *AfterSaleOrderSyncResponse) MetaData() map[string]any { return make(map[string]any) }
