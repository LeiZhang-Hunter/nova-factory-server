package model

import (
	"encoding/json"
	"nova-factory-server/app/utils/observer/integration/result"
)

// OrderStatusSyncResult 单笔订单状态同步结果，实现 result.OrderStatusSyncResult。
type OrderStatusSyncResult struct {
	Tid     string `json:"tid"`
	Message string `json:"message"`
}

func (r *OrderStatusSyncResult) GetTid() string     { return r.Tid }
func (r *OrderStatusSyncResult) GetMessage() string { return r.Message }
func (r *OrderStatusSyncResult) Ptr() any           { return r }
func (r *OrderStatusSyncResult) RawStr() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func (r *OrderStatusSyncResult) MetaData() map[string]any { return make(map[string]any) }

// OrderStatusSyncResponse 订单状态同步响应，实现 result.OrderStatusSyncResponse。
type OrderStatusSyncResponse struct {
	Code    int64                    `json:"code"`
	Message string                   `json:"message"`
	Orders  []*OrderStatusSyncResult `json:"orders"`
}

func (r *OrderStatusSyncResponse) GetCode() int64     { return r.Code }
func (r *OrderStatusSyncResponse) GetMessage() string { return r.Message }
func (r *OrderStatusSyncResponse) GetOrders() []result.OrderStatusSyncResult {
	out := make([]result.OrderStatusSyncResult, len(r.Orders))
	for i, v := range r.Orders {
		out[i] = v
	}
	return out
}
func (r *OrderStatusSyncResponse) Ptr() any { return r }
func (r *OrderStatusSyncResponse) RawStr() (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func (r *OrderStatusSyncResponse) MetaData() map[string]any { return make(map[string]any) }
