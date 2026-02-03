package deviceMonitorServiceImpl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	controlService "github.com/novawatcher-io/nova-factory-payload/control/v1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/deviceMonitor/deviceMonitorService"
	"nova-factory-server/app/business/metric/device/metricDao"
	"nova-factory-server/app/business/metric/device/metricModels"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/datasource/cache"
	"sync"
	"time"
)

type DeviceControlServiceImpl struct {
	manager         *ControlManagerServiceImpl
	agentDao        daemonizeDao.IotAgentDao
	controlLogDao   metricDao.IControlLogDao
	cache           cache.Cache
	waitMap         sync.Map
	dao             deviceDao.IDeviceDao
	deviceConfigDao deviceDao.ISysModbusDeviceConfigDataDao
}

func NewDeviceControlServiceImpl(agentDao daemonizeDao.IotAgentDao, cache cache.Cache,
	controlLogDao metricDao.IControlLogDao,
	dao deviceDao.IDeviceDao,
	deviceConfigDao deviceDao.ISysModbusDeviceConfigDataDao) deviceMonitorService.DeviceControlService {
	return &DeviceControlServiceImpl{
		manager:         NewControlManagerServiceImpl(),
		agentDao:        agentDao,
		cache:           cache,
		controlLogDao:   controlLogDao,
		dao:             dao,
		deviceConfigDao: deviceConfigDao,
	}
}

func (d *DeviceControlServiceImpl) Control(ctx context.Context, request *controlService.ControlRequest) (*controlService.ControlResponse, error) {
	if request == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "request is nil")
	}

	key := fmt.Sprintf(device.DEVICE_CONTROL_KEY, request.DeviceId, request.DataId)
	d.cache.Del(context.Background(), key)

	waitChan, ok := d.waitMap.Load(request.RequestId)
	if !ok {
		return nil, nil
	}

	reqChan, ok := waitChan.(chan *controlService.ControlRequest)
	if !ok {
		return nil, nil
	}
	reqChan <- request
	gctx := &gin.Context{}
	// 读取设备信息
	deviceInfo, err := d.dao.GetById(gctx, int64(request.DeviceId))
	if err != nil {
		zap.L().Error("get device error", zap.Error(err))
		return nil, errors.New("读取设备信息错误")
	}
	if deviceInfo == nil {
		return nil, errors.New("设备不存在")
	}

	var deviceName string
	if deviceInfo.Name != nil {
		deviceName = *deviceInfo.Name
	}

	dataInfo, err := d.deviceConfigDao.GetById(gctx, request.DataId)
	if err != nil {
		zap.L().Error("get device data error", zap.Error(err))
		return nil, errors.New("读取设备模板错误")
	}

	if dataInfo == nil {
		zap.L().Error("get device data error", zap.Error(err))
		return nil, errors.New("设备模板信息不存在")
	}

	content, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New("下发内容编码失败")
	}

	serieId := time.Now().UnixNano()
	var now *gtime.Time = gtime.Now()
	err = d.controlLogDao.Export(context.Background(), []*metricModels.NovaControlLog{
		{
			DeviceId:   request.DeviceId,
			DeviceName: deviceName,
			DataId:     request.DataId,
			DataName:   dataInfo.Name,
			Message: fmt.Sprintf("控制<%s><%s>指令 value %s 下发回执;返回错误信息:%s", deviceName, dataInfo.Name,
				string(content), request.Value.GetStringValue()),
			Type:          "manual",
			SeriesId:      uint64(serieId),
			Attributes:    make(map[string]string),
			StartTimeUnix: now,
			TimeUnix:      now,
		},
	})
	if err != nil {
		zap.L().Error("Export error", zap.Error(err))
		return nil, err
	}
	return nil, nil
}

func (d *DeviceControlServiceImpl) Register(context.Context, *controlService.RegisterReq) (*controlService.RegisterRes, error) {
	return nil, nil
}
func (d *DeviceControlServiceImpl) Unregister(context.Context, *controlService.UnregisterReq) (*controlService.UnregisterRes, error) {
	return nil, nil
}
func (d *DeviceControlServiceImpl) Heartbeat(context.Context, *controlService.HeartbeatReq) (*controlService.HeartbeatRes, error) {
	return nil, nil
}
func (d *DeviceControlServiceImpl) Operate(req *controlService.OperateReq, stream controlService.ControlService_OperateServer) error {
	//return nil
	clientId := req.GetObjectId()
	if clientId == 0 {
		zap.L().Error("AgentOperate error: agent id is 0")
		return gerror.NewCode(gcode.CodeInvalidParameter, "agent id cannot be empty")
	}

	g.Log().Debugf(stream.Context(), "agent add client stream, id: %v", clientId)
	d.manager.AddClient(clientId, stream)
	defer d.manager.DeleteClient(clientId)
	for {
		if d.manager.IsStopped() {
			g.Log().Debugf(stream.Context(), "server stopped, release client stream, id: %v", clientId)
			return nil
		}

		select {
		case <-stream.Context().Done():
			{
				return nil
			}
		}
	}
}

func (d *DeviceControlServiceImpl) BroadcastControl(ctx context.Context, req *controlService.ControlRequest) (*controlService.ControlResponse, error) {
	if req == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "req is nil")
	}

	streamClient := d.manager.getClient(req.AgentId)
	if streamClient == nil {
		return &controlService.ControlResponse{
			RequestId: req.RequestId,
			Code:      404,
		}, nil
	}

	// redis 变成发送中
	//key := fmt.Sprintf(device.DEVICE_CONTROL_KEY, req.DeviceId, req.DataId)
	//res := d.cache.SetNX(ctx, key, 1, 15*time.Second)
	//if !res {
	//	return &controlService.ControlResponse{
	//		RequestId: req.RequestId,
	//		Code:      201,
	//	}, nil
	//}

	err := streamClient.Send(&controlService.OperateRes{
		Request: &controlService.ControlRequest{
			RequestId: req.RequestId,
			DeviceId:  req.DeviceId,
			AgentId:   req.AgentId,
			DataId:    req.DataId,
			Timestamp: timestamppb.Now(),
			Value:     req.Value,
		},
	})
	if err != nil {
		zap.L().Error("deviceMonitorServiceImpl.OperateRes", zap.Error(err))
		return nil, err
	}

	// 等待下游返回结果
	var waitChan chan *controlService.ControlRequest = make(chan *controlService.ControlRequest)
	d.waitMap.Store(req.RequestId, waitChan)
	select {
	case request := <-waitChan:
		{
			d.waitMap.Delete(req.RequestId)
			if request.GetValue().GetStringValue() == "" {
				return &controlService.ControlResponse{
					RequestId: request.RequestId,
					Code:      0,
				}, nil
			} else {
				return &controlService.ControlResponse{
					RequestId: request.RequestId,
					Message:   request.GetValue().GetStringValue(),
					Code:      500,
				}, nil
			}
		}
	case <-time.After(5 * time.Second):
		{
			d.waitMap.Delete(req.RequestId)
			return &controlService.ControlResponse{
				RequestId: req.RequestId,
				Code:      503,
				Message:   "控制指令下发超时",
			}, nil
		}

	}
}
