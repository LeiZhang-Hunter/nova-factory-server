package craftRouteServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCraftRouteServiceImpl, NewICraftProcessServiceImpl, NewICraftProcessContextServiceImpl,
	NewIProcessRouteServiceImpl,
	NewIScheduleServiceImpl)
