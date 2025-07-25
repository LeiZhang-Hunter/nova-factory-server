package craftRouteDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCraftRouteDaoImpl, NewIProcessDaoImpl, NewProcessContextDaoImpl,
	NewIProcessRouteDaoImpl, NewISysProRouteProductDaoImpl, NewSysProRouteProductBomDaoImpl,
	NewISysCraftRouteConfigDaoImpl, NewWorkOrderDaoImpl, NewISysProTaskDaoImpl, NewIScheduleDaoImpl, NewIScheduleMapDaoImpl)
