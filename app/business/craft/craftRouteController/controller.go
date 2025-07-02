package craftRouteController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCraft, NewProcess, NewProcessContext, NewSysProRouteProcess,
	NewRouteProduct, NewRouteProductBom, NewWorkOrder, NewTask,
	wire.Struct(new(CraftRoute), "*"))

type CraftRoute struct {
	CraftRoute      *Craft
	Process         *Process
	ProcessContext  *ProcessContext
	RouteProcess    *RouteProcess
	RouteProduct    *RouteProduct
	RouteProductBom *RouteProductBom
	WorkOrder       *WorkOrder
	Task            *Task
}
