package settingServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAgentConfigService)
