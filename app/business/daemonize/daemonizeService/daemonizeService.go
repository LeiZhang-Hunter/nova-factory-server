package daemonizeService

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/daemonize/grpc/v1"
)

type DaemonizeService interface {
	AgentRegister(ctx context.Context, req *v1.AgentRegisterReq) (res *v1.AgentRegisterRes, err error)
}
