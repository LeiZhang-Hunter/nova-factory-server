package shopcontroller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewUser, wire.Struct(new(Controller), "*"))

type Controller struct {
	User *User
}
