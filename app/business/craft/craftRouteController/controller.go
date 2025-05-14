package craftRouteController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCraft, NewProcess, NewProcessContext, wire.Struct(new(CraftRoute), "*"))

type CraftRoute struct {
	CraftRoute     *Craft
	Process        *Process
	ProcessContext *ProcessContext
}
