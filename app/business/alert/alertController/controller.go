package alertController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAlert, wire.Struct(new(Controller), "*"))

type Controller struct {
	Alert *Alert
}
