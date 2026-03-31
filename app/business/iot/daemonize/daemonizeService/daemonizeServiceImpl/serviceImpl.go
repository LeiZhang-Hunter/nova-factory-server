package daemonizeServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDaemonizeServiceImpl, NewManagerServiceImpl,
	NewIotAgentServiceImpl, NewIGatewayConfigServiceImpl, NewIotAgentConfigServiceImpl)
