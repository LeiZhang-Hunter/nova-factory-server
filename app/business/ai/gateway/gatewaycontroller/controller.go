package gatewaycontroller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIGateway, NewAgent, wire.Struct(new(Controller), "*"))

type Controller struct {
	AIGateway *AIGateway
	Agent     *Agent
}
