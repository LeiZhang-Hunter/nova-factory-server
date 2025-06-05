package daemonizeController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDaemonize, wire.Struct(new(DaemonizeServer), "*"))

type DaemonizeServer struct {
	Daemonize *Daemonize
}
