package gatewaycontroller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIGateway, NewAgent, NewSkills, NewMCPServer,
	NewConversations, NewMessage, wire.Struct(new(Controller), "*"))

type Controller struct {
	AIGateway     *AIGateway
	Agent         *Agent
	Skills        *Skills
	MCPServer     *MCPServer
	Conversations *Conversations
	Message       *Message
}
