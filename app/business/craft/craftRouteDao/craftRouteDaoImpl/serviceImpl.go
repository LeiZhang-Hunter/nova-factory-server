package craftRouteDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCraftRouteDaoImpl, NewIProcessDaoImpl, NewProcessContextDaoImpl,
	NewIProcessRouteDaoImpl, NewISysCraftRouteConfigDaoImpl, NewIScheduleDaoImpl, NewIScheduleMapDaoImpl)
