package gatewayserviceimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIGatewayService, NewIAiConversationServiceImpl, NewInstalledSkillService, NewMCPServerService)
