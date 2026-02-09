package configurationController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewConfiguration, wire.Struct(new(Controller), "*"))

type Controller struct {
	Configuration *Configuration
}
