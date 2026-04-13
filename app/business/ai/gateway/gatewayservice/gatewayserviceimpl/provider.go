package gatewayserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIAgentService, NewAIGatewayService, NewIAiConversationServiceImpl, NewInstalledSkillService, NewMCPServerService)
