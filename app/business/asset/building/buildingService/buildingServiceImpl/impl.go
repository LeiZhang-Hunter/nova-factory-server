package buildingServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewBuildingServiceImpl, NewFloorServiceImpl)
