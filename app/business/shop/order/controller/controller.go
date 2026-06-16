package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewOrder, NewAutoCancel, wire.Struct(new(Controller), "*"))

type Controller struct {
	Order      *Order
	AutoCancel *AutoCancel
}
