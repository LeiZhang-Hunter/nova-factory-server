package gatewaydaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIAgentDao, NewAIAgentOrchestrationDao, NewAISubAgentDao, NewAIGatewayDao, NewAIAgentMessageDao, NewInstalledSkillDao, NewMCPServerDao, NewAIAgentConfigPublishHistoryDao, NewAgentConfigKeyDao)
