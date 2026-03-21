package settingController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAgentConfig, wire.Struct(new(Controller), "*"))

type Controller struct {
	AgentConfig *AgentConfig
}
