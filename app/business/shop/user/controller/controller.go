package shopcontroller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewAddress, NewCart, NewOrder, NewUser, wire.Struct(new(Controller), "*"))

type Controller struct {
	Address *Address
	Cart    *Cart
	Order   *Order
	User    *User
}
