package craftRouteController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCraft, NewProcess, NewProcessContext, NewSysProRouteProcess, wire.Struct(new(CraftRoute), "*"))

type CraftRoute struct {
	CraftRoute     *Craft
	Process        *Process
	ProcessContext *ProcessContext
	RouteProcess   *RouteProcess
}
