package daemonizeDao

import (
	"context"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
)

type IotAgentProcess interface {
	RecordHeardBeat(ctx context.Context, objectId uint64, processes []*daemonizeModels.SysIotAgentProcess) error
}
