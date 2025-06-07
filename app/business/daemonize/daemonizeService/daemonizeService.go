package daemonizeService

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/daemonize/grpc/v1"
)

type DaemonizeService interface {
	AgentRegister(ctx context.Context, req *v1.AgentRegisterReq) (res *v1.AgentRegisterRes, err error)
	AgentGetConfig(ctx context.Context, req *v1.AgentGetConfigReq) (*v1.AgentGetConfigRes, error)
	AgentHeartbeat(ctx context.Context, req *v1.AgentHeartbeatReq) (res *v1.AgentHeartbeatRes, err error)
	AgentOperate(ctx context.Context, req *v1.AgentOperateReq, stream v1.AgentControllerService_AgentOperateServer) (err error)
	AgentOperateProcess(ctx context.Context, cmd v1.AgentCmd, processOperateInfoList []*v1.ProcessOperateInfo)
	BroadcastAgentOperateProcess(ctx context.Context, cmd v1.AgentCmd, processOperateInfoList []*v1.ProcessOperateInfo) error
}
