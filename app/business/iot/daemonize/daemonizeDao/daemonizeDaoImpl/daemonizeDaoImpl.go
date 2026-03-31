package daemonizeDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIotAgentDaoImpl, NewIotAgentConfigDaoImpl, NewIotAgentProcessDaoImpl)
