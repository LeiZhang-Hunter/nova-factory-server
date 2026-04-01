package gatewaycontroller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIGateway, wire.Struct(new(Controller), "*"))

type Controller struct {
	AIGateway *AIGateway
}
