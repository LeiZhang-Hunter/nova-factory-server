package settingController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAgentConfig, NewIntegrationConfig, wire.Struct(new(Controller), "*"))

type Controller struct {
	AgentConfig       *AgentConfig
	IntegrationConfig *IntegrationConfig
}
