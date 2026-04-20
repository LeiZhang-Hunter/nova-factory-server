package shopcontroller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewAddress, NewCart, NewUser, wire.Struct(new(Controller), "*"))

type Controller struct {
	Address *Address
	Cart    *Cart
	User    *User
}
