package gatewaycontroller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIGateway, NewAgent, NewSkills, NewMCPServer,
	NewAgentOrchestration, NewSubAgent, NewConversations, NewMessage, NewAgentConfigPublishHistory,
	NewAgentConfigKey, wire.Struct(new(Controller), "*"))

type Controller struct {
	AIGateway            *AIGateway
	Agent                *Agent
	Orchestration        *AgentOrchestration
	SubAgent             *SubAgent
	Skills               *Skills
	MCPServer            *MCPServer
	Conversations        *Conversations
	Message              *Message
	ConfigPublishHistory *AgentConfigPublishHistory
	AgentConfigKey       *AgentConfigKey
}
