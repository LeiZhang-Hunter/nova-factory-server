package daemonizeDao

import (
	"context"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
)

type IotAgentDao interface {
	GetByObjectId(ctx context.Context, objectId uint64) (agent *daemonizeModels.SysIotAgent, err error)
	UpdateHeartBeat(ctx context.Context, data *daemonizeModels.SysIotAgent) error
	Update(ctx context.Context, data *daemonizeModels.SysIotAgent) error
}
