package dashboardDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDashboardDaoImpl, NewDashboardDataDaoImpl)
