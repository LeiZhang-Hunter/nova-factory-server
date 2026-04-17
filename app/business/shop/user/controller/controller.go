package shopcontroller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAddress, NewUser, wire.Struct(new(Controller), "*"))

type Controller struct {
	Address *Address
	User    *User
}
