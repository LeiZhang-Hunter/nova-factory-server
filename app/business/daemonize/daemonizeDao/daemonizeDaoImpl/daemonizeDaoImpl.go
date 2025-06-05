package daemonizeDao

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIotAgentDaoImpl)
