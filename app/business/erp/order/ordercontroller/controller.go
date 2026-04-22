package ordercontroller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewOrder, NewOrderAudit, wire.Struct(new(Controller), "*"))

type Controller struct {
	Order      *Order
	OrderAudit *OrderAudit
}
