package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCombination, NewPink, wire.Struct(new(Controller), "*"))

type Controller struct {
	Combination *Combination
	Pink        *Pink
}
