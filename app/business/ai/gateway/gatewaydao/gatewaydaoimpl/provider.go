package gatewaydaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIAgentDao, NewAISubAgentDao, NewAIGatewayDao, NewAIAgentMessageDao, NewInstalledSkillDao, NewMCPServerDao)
