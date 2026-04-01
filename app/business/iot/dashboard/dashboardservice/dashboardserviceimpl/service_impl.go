package dashboardserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDashboardServiceImpl, NewDashboardDataServiceImpl)
