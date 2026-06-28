package controller

import "github.com/google/wire"

type Controller struct {
	Tracking *Tracking
}

var ProviderSet = wire.NewSet(NewTracking, wire.Struct(new(Controller), "*"))
