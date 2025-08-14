package buildingController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewBuilding,
	wire.Struct(new(Building), "*"))

type Controller struct {
	Building *Building
}
