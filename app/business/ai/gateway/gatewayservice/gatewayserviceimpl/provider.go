package gatewayserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIAgentService, NewAIAgentOrchestrationService, NewAISubAgentService, NewAIGatewayService, NewIAiConversationServiceImpl,
	NewInstalledSkillService, NewMCPServerService, NewIAiMessageServiceImpl)
