package gatewaycontroller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIGateway, NewConversations, wire.Struct(new(Controller), "*"))

type Controller struct {
	AIGateway     *AIGateway
	Conversations *Conversations
}
