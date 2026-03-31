package orderController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewOrder, wire.Struct(new(Controller), "*"))

type Controller struct {
	Order *Order
}
