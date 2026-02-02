package deviceMonitorController

import (
	"context"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	controlService "github.com/novawatcher-io/nova-factory-payload/control/v1"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
)

type DeviceControl struct {
	service deviceMonitorService.DeviceControlService
	controlService.UnimplementedControlServiceServer
}

func NewDeviceControl(service deviceMonitorService.DeviceControlService) *DeviceControl {
	return &DeviceControl{
		service: service,
	}
}

func (d *DeviceControl) PrivateRoutes(router *grpcx.GrpcServer) {
	controlService.RegisterControlServiceServer(router.Server, d)
}

// Control Bidirectional stream for real-time control and feedback.
func (d *DeviceControl) Control(ctx context.Context, request *controlService.ControlRequest) (*controlService.ControlResponse, error) {
	return d.service.Control(ctx, request)
}

// BroadcastControl Unary call for single control operations.
func (d *DeviceControl) BroadcastControl(ctx context.Context, req *controlService.ControlRequest) (*controlService.ControlResponse, error) {
	return d.service.BroadcastControl(ctx, req)
}

// Register 注册 agent
func (d *DeviceControl) Register(context.Context, *controlService.RegisterReq) (*controlService.RegisterRes, error) {
	return nil, nil
}

// Unregister 注销 agent
func (d *DeviceControl) Unregister(context.Context, *controlService.UnregisterReq) (*controlService.UnregisterRes, error) {
	return nil, nil
}

// Heartbeat agent 心跳
func (d *DeviceControl) Heartbeat(context.Context, *controlService.HeartbeatReq) (*controlService.HeartbeatRes, error) {
	return nil, nil
}

// Operate 操作agent stream长连接
func (d *DeviceControl) Operate(req *controlService.OperateReq, stream controlService.ControlService_OperateServer) error {
	return d.service.Operate(req, stream)
}
