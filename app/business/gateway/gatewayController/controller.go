package gatewayController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewProtocol, wire.Struct(new(GatewayController), "*"))

type GatewayController struct {
	Protocol *Protocol
}
