package deviceMonitorController

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	daemonizeService "nova-factory-server/app/business/iot/daemonize/daemonizeService"
	"nova-factory-server/app/constant/agent"
	"nova-factory-server/app/utils/uuid"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	client "github.com/novawatcher-io/nova-factory-payload/camera/v1"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
)

type cameraJsonCodec struct{}

func (cameraJsonCodec) Name() string {
	return "json"
}

func (cameraJsonCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (cameraJsonCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

var cameraCodecRegisterOnce sync.Once

// CameraGrpc 提供摄像头信令订阅、广播与令牌回传能力。
type CameraGrpc struct {
	client.UnimplementedCameraServiceServer
	mu           sync.RWMutex
	agentService daemonizeService.IotAgentService
	manager      *CameraSubscribeManager
	waitMap      sync.Map
	tokenWaiters map[string]chan client.TokenAck
	lastAck      map[string]client.TokenAck
}

// NewCameraGrpc 创建 CameraGrpc 实例。
func NewCameraGrpc(agentService daemonizeService.IotAgentService) *CameraGrpc {
	return &CameraGrpc{
		agentService: agentService,
		manager:      NewCameraSubscribeManager(),
		tokenWaiters: make(map[string]chan client.TokenAck),
		lastAck:      make(map[string]client.TokenAck),
	}
}

// PrivateRoutes 注册 CameraService gRPC 服务。
func (c *CameraGrpc) PrivateRoutes(router *grpcx.GrpcServer) {
	cameraCodecRegisterOnce.Do(func() {
		encoding.RegisterCodec(cameraJsonCodec{})
	})
	client.RegisterCameraServiceServer(router.Server, c)
	info := router.Server.GetServiceInfo()
	fmt.Println(info)
	return
}

// Report 兼容上报接口，当前返回成功占位。
func (c *CameraGrpc) Report(_ context.Context, _ *client.CameraData) (*client.CameraResponse, error) {
	return &client.CameraResponse{Code: 0}, nil
}

// WebrtcSubscribe 建立下游订阅连接并注册到订阅管理器。
func (c *CameraGrpc) WebrtcSubscribe(req *client.SubscribeRequest, stream client.CameraService_WebrtcSubscribeServer) error {
	node := strings.TrimSpace(req.Node)
	if node == "" {
		node = "default"
	}
	c.manager.AddClient(node, stream)
	defer c.manager.DeleteClient(node, stream)
	<-stream.Context().Done()
	return nil
}

// WebrtcSendToken 接收下游回传令牌并唤醒等待中的请求。
func (c *CameraGrpc) WebrtcSendToken(_ context.Context, ack *client.TokenAck) (*client.SendTokenReply, error) {
	ch, ok := c.waitMap.Load(ack.RequestId)
	if !ok {
		return &client.SendTokenReply{}, errors.New("waitMap load fail")
	}
	if ch == nil {
		return &client.SendTokenReply{}, errors.New("waitMap load fail")
	}
	defer func() {
		ch.(chan *client.TokenAck) <- ack
	}()
	if ack.Code != 0 {
		return &client.SendTokenReply{}, errors.New(ack.Errmsg)
	}

	if ack == nil {
		return &client.SendTokenReply{}, errors.New("ack is nil")
	}
	if ack.DeviceId == "" {
		return &client.SendTokenReply{}, errors.New("ack.DeviceId is empty")
	}

	return &client.SendTokenReply{}, nil
}

// WebrtcBroadcast 将协商请求投递到目标节点订阅流。
func (c *CameraGrpc) WebrtcBroadcast(_ context.Context, req *client.WebrtcBroadcastRequest) (*client.WebrtcBroadcastReply, error) {
	if req == nil {
		return &client.WebrtcBroadcastReply{Code: -1, Msg: "request is nil"}, nil
	}
	message := client.SubscribeMessage{
		RequestId: req.GetRequestId(),
		DeviceId:  req.GetDeviceId(),
		ChannelId: req.GetChannelId(),
		Sdp64:     req.GetSdp64(),
	}
	targetNode := strings.TrimSpace(req.GetTargetNode())
	if targetNode == "" {
		return &client.WebrtcBroadcastReply{
			Code: -1,
			Msg:  "请选择",
		}, nil
	}
	var waitChan chan *client.TokenAck = make(chan *client.TokenAck)
	c.waitMap.Store(message.RequestId, waitChan)
	defer c.waitMap.Delete(message.RequestId)
	delivered, err := c.manager.PublishToNode(targetNode, &message)
	if err != nil {
		return &client.WebrtcBroadcastReply{Code: -1, Msg: err.Error()}, nil
	}
	ack := <-waitChan
	if ack.Code != 0 {
		return &client.WebrtcBroadcastReply{Code: -1, Msg: ack.Errmsg}, errors.New(ack.Errmsg)
	}
	return &client.WebrtcBroadcastReply{
		Code:           0,
		Msg:            "ok",
		DeliveredCount: delivered,
		DeliveredNodes: []string{targetNode},
		Sdp64:          ack.Sdp64,
	}, nil
}

// PublishStart 对外保持兼容，内部统一走广播协商流程。
func (c *CameraGrpc) PublishStart(node string, message *client.SubscribeMessage, timeout time.Duration) (*client.TokenAck, error) {
	return c.PublishStartByBroadcast(node, message, timeout)
}

// PublishStartByBroadcast 广播 offer 并等待下游 token 回传。
func (c *CameraGrpc) PublishStartByBroadcast(node string, message *client.SubscribeMessage, timeout time.Duration) (*client.TokenAck, error) {
	if message.DeviceId == "" {
		return nil, errors.New("device_id is required")
	}
	node = strings.TrimSpace(node)
	requestId := uuid.MakeUuid()
	broadcastReq := &client.WebrtcBroadcastRequest{
		RequestId:  requestId,
		Source:     "http_camera_offer",
		TargetNode: node,
		DeviceId:   message.DeviceId,
		ChannelId:  message.ChannelId,
		Sdp64:      message.Sdp64,
	}
	broadcastReply, err := c.broadcastByGrpcClient(broadcastReq)
	if err != nil {
		return nil, err
	}
	if broadcastReply == nil || broadcastReply.GetCode() != 0 || broadcastReply.GetDeliveredCount() == 0 {
		return nil, errors.New("broadcast to camera subscribers failed")
	}

	return &client.TokenAck{
		ChannelId: message.ChannelId,
		DeviceId:  message.DeviceId,
		Sdp64:     broadcastReply.Sdp64,
	}, nil
}

// broadcastByGrpcClient 通过集群地址列表并发调用广播接口。
func (c *CameraGrpc) broadcastByGrpcClient(req *client.WebrtcBroadcastRequest) (*client.WebrtcBroadcastReply, error) {
	addressList := viper.GetStringSlice("daemonize.server_list")
	if len(addressList) == 0 {
		return nil, errors.New("daemonize.server_list is empty")
	}
	if c.agentService == nil {
		return nil, errors.New("agent service is nil")
	}
	objectID, err := strconv.ParseUint(strings.TrimSpace(req.GetTargetNode()), 10, 64)
	if err != nil || objectID == 0 {
		return nil, errors.New("target node invalid")
	}
	info, err := c.agentService.GetByObjectId(context.Background(), objectID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, errors.New("agent not found")
	}
	var wg sync.WaitGroup
	wg.Add(len(addressList))
	var mu sync.Mutex
	var replayErr error
	var successReply *client.WebrtcBroadcastReply
	for _, address := range addressList {
		targetAddress := address
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			grpcCtx := metadata.AppendToOutgoingContext(ctx,
				agent.USERNAME, info.Username,
				agent.PASSWORD, info.Password,
				agent.GATEWAYID, strconv.FormatUint(info.ObjectID, 10),
			)
			var reply *client.WebrtcBroadcastReply
			rpcClient := client.NewCameraServiceClient(grpcx.Client.MustNewGrpcClientConn(targetAddress))
			reply, replayErr = rpcClient.WebrtcBroadcast(grpcCtx, req)
			if replayErr != nil {
				zap.L().Error("广播失败", zap.Error(err))
				return
			}
			if reply == nil || reply.GetCode() != 0 || reply.GetDeliveredCount() == 0 {
				return
			}
			mu.Lock()
			if successReply == nil {
				successReply = reply
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	if successReply != nil {
		return successReply, nil
	}
	if replayErr != nil {
		return nil, replayErr
	}
	return nil, errors.New("broadcast to camera subscribers failed")
}

// GetLastToken 读取最近一次会话回传的 token 数据。
func (c *CameraGrpc) GetLastToken(deviceID, channelID string) (*client.TokenAck, bool) {
	key := c.sessionKey(deviceID, channelID)
	c.mu.RLock()
	ack, ok := c.lastAck[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return &ack, true
}

// sessionKey 生成设备与通道维度的会话键。
func (c *CameraGrpc) sessionKey(deviceID, channelID string) string {
	return deviceID + ":" + channelID
}
