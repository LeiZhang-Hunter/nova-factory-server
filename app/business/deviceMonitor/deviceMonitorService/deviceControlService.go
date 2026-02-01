package deviceMonitorService

import (
	"context"
	controlService "github.com/novawatcher-io/nova-factory-payload/control/v1"
)

type DeviceControlService interface {
	// Control 控制结果反馈，删除redis标记
	Control(ctx context.Context, request *controlService.ControlRequest) (*controlService.ControlResponse, error)
	// BroadcastControl 广播agent，下发控制指令
	BroadcastControl(context.Context, *controlService.ControlRequest) (*controlService.ControlResponse, error)
	Register(context.Context, *controlService.RegisterReq) (*controlService.RegisterRes, error)
	Unregister(context.Context, *controlService.UnregisterReq) (*controlService.UnregisterRes, error)
	Heartbeat(context.Context, *controlService.HeartbeatReq) (*controlService.HeartbeatRes, error)
	// Operate stream长连接
	Operate(*controlService.OperateReq, controlService.ControlService_OperateServer) error
}
