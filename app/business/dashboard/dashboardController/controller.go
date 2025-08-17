package dashboardController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDashboard, wire.Struct(new(Controller), "*"))

type Controller struct {
	Dashboard *Dashboard
}
