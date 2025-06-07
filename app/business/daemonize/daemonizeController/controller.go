package daemonizeController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDaemonize, NewIotAgentController, wire.Struct(new(DaemonizeServer), "*"))

type DaemonizeServer struct {
	Daemonize *Daemonize
	IotAgent  *IotAgent
}
