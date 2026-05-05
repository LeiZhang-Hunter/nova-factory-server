package shopcontroller

import "github.com/google/wire"

type Controller struct {
	Address *Address
	Cart    *Cart
	User    *User
}

var ProviderSet = wire.NewSet(NewAddress, NewCart, NewUser, wire.Struct(new(Controller), "*"))
