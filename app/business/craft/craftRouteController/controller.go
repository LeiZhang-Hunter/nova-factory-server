package craftRouteController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCraft, wire.Struct(new(CraftRoute), "*"))

type CraftRoute struct {
	CraftRoute *Craft
}
