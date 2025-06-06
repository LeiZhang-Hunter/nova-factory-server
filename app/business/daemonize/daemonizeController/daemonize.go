package daemonizeController

import (
	"context"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	v1 "github.com/novawatcher-io/nova-factory-payload/daemonize/grpc/v1"
	"nova-factory-server/app/business/daemonize/daemonizeService"
)

type Daemonize struct {
	service daemonizeService.DaemonizeService
	v1.UnimplementedAgentControllerServiceServer
}

func NewDaemonize(service daemonizeService.DaemonizeService) *Daemonize {
	return &Daemonize{
		service: service,
	}
}

func (d *Daemonize) PrivateRoutes(router *grpcx.GrpcServer) {
	v1.RegisterAgentControllerServiceServer(router.Server, d)
}

// AgentRegister 注册 agent
func (d *Daemonize) AgentRegister(ctx context.Context, req *v1.AgentRegisterReq) (*v1.AgentRegisterRes, error) {
	return d.service.AgentRegister(ctx, req)
}

// AgentUnregister 注销 agent
func (d *Daemonize) AgentUnregister(context.Context, *v1.AgentUnregisterReq) (*v1.AgentUnregisterRes, error) {
	return nil, nil
}

// AgentHeartbeat agent 心跳
func (d *Daemonize) AgentHeartbeat(c context.Context, req *v1.AgentHeartbeatReq) (*v1.AgentHeartbeatRes, error) {
	return d.service.AgentHeartbeat(c, req)
}

// AgentOperate 操作agent stream长连接
func (d *Daemonize) AgentOperate(req *v1.AgentOperateReq, stream v1.AgentControllerService_AgentOperateServer) error {
	return d.service.AgentOperate(stream.Context(), req, stream)
}

// AgentGetConfig 获取配置
func (d *Daemonize) AgentGetConfig(c context.Context, req *v1.AgentGetConfigReq) (*v1.AgentGetConfigRes, error) {
	return d.service.AgentGetConfig(c, req)
}

// AgentOperateProcess 操作agent的子进程 广播使用
func (d *Daemonize) AgentOperateProcess(ctx context.Context, req *v1.AgentOperateProcessReq) (*v1.AgentOperateProcessRes, error) {
	d.service.AgentOperateProcess(ctx, req.Cmd, req.OperateInfoList)
	return &v1.AgentOperateProcessRes{}, nil
}
