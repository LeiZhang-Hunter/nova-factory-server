package daemonizedao

import (
	"context"
	"nova-factory-server/app/business/iot/daemonize/daemonizemodels"
)

type IotAgentProcess interface {
	RecordHeardBeat(ctx context.Context, objectId uint64, processes []*daemonizemodels.SysIotAgentProcess) error
	GetHeardBeatInfo(ctx context.Context, objectIds []uint64) map[uint64][]*daemonizemodels.SysIotAgentProcess
}
