package dashboarddaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDashboardDaoImpl, NewDashboardDataDaoImpl)
