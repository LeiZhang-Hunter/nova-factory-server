package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewElectric, wire.Struct(new(System), "*"))

type System struct {
	Electric *Electric
}
