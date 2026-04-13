package gatewaydaoimpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAIAgentDao, NewAIGatewayDao, NewAIAgentMessageDao, NewInstalledSkillDao, NewMCPServerDao)
