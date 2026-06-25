package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewAccount,
	wire.Struct(new(Controller), "*"),
)

type Controller struct {
	Account *Account
}
