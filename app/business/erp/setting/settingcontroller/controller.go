package settingcontroller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAgentConfig, NewIntegrationConfig, NewLogistics, wire.Struct(new(Controller), "*"))

type Controller struct {
	AgentConfig       *AgentConfig
	IntegrationConfig *IntegrationConfig
	Logistics         *Logistics
}
