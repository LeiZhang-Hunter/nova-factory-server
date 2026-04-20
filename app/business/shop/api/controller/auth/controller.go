package auth

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAuth, wire.Struct(new(Controller), "*"))

type Controller struct {
	Auth *Auth
}
