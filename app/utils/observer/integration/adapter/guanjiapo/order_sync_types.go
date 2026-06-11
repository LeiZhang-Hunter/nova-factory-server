package guanjiapo

import (
	"encoding/json"
	"nova-factory-server/app/utils/observer/integration/result"
)

type OrderSyncResult struct {
	Tid      string `json:"tid"`
	BillCode string `json:"billcode"`
	Message  string `json:"message"`
}

type OrderSyncResponse struct {
	Code    int64                    `json:"code"`
	Message string                   `json:"message"`
	Orders  []result.OrderSyncResult `json:"orders"`
}

func (o *OrderSyncResponse) GetCode() int64 {
	return o.Code
}
func (o *OrderSyncResponse) GetMessage() string {
	return o.Message
}
func (o *OrderSyncResponse) GetOrders() []result.OrderSyncResult {
	return o.Orders
}

func (o *OrderSyncResponse) Ptr() any {
	return o
}
func (o *OrderSyncResponse) RawStr() (string, error) {
	marshal, err := json.Marshal(o)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
func (o *OrderSyncResponse) MetaData() map[string]any {
	return make(map[string]any)
}
