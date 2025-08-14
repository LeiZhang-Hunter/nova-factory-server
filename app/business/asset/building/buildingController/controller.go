package buildingController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewBuilding,
	wire.Struct(new(Controller), "*"))

type Controller struct {
	Building *Building
}
