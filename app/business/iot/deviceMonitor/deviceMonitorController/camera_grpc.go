package deviceMonitorController

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
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

type subscribeRequest struct {
	Node string `json:"node"`
}

type subscribeMessage struct {
	DeviceId  string `json:"device_id"`
	ChannelId string `json:"channel_id"`
	SDP64     string `json:"sdp64"`
}

type tokenAck struct {
	DeviceId  string `json:"device_id"`
	ChannelId string `json:"channel_id"`
	Token     string `json:"token"`
	PlayURL   string `json:"play_url"`
	WhepURL   string `json:"whep_url"`
	SDP64     string `json:"sdp64"`
}

type cameraAckResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type webRTCSubscribeServer interface {
	Send(*subscribeMessage) error
	grpc.ServerStream
}

type webRTCServiceServer interface {
	Subscribe(*subscribeRequest, webRTCSubscribeServer) error
	SendToken(context.Context, *tokenAck) (*cameraAckResponse, error)
}

type webRTCSubscribeServerWrapper struct {
	grpc.ServerStream
}

func (x *webRTCSubscribeServerWrapper) Send(m *subscribeMessage) error {
	return x.ServerStream.SendMsg(m)
}

type CameraGrpc struct {
	subscriberSeq uint64
	mu            sync.RWMutex
	subscribers   map[string]map[uint64]chan subscribeMessage
	tokenWaiters  map[string]chan tokenAck
	lastAck       map[string]tokenAck
}

func NewCameraGrpc() *CameraGrpc {
	return &CameraGrpc{
		subscribers:  make(map[string]map[uint64]chan subscribeMessage),
		tokenWaiters: make(map[string]chan tokenAck),
		lastAck:      make(map[string]tokenAck),
	}
}

func (c *CameraGrpc) PrivateRoutes(router *grpcx.GrpcServer) {
	cameraCodecRegisterOnce.Do(func() {
		encoding.RegisterCodec(cameraJsonCodec{})
	})
	router.Server.RegisterService(&_WebRTCService_serviceDesc, c)
	info := router.Server.GetServiceInfo()
	fmt.Println(info)
	return
}

func (c *CameraGrpc) Subscribe(req *subscribeRequest, stream webRTCSubscribeServer) error {
	node := strings.TrimSpace(req.Node)
	if node == "" {
		node = "default"
	}
	ch := make(chan subscribeMessage, 16)
	subID := atomic.AddUint64(&c.subscriberSeq, 1)
	c.addSubscriber(node, subID, ch)
	defer c.removeSubscriber(node, subID)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg := <-ch:
			if err := stream.Send(&msg); err != nil {
				return err
			}
		}
	}
}

func (c *CameraGrpc) SendToken(_ context.Context, ack *tokenAck) (*cameraAckResponse, error) {
	if ack == nil {
		return &cameraAckResponse{Code: -1, Message: "ack is nil"}, nil
	}
	if ack.DeviceId == "" {
		return &cameraAckResponse{Code: -1, Message: "device_id is required"}, nil
	}
	key := c.sessionKey(ack.DeviceId, ack.ChannelId)
	c.mu.Lock()
	c.lastAck[key] = *ack
	waiter, ok := c.tokenWaiters[key]
	if ok {
		delete(c.tokenWaiters, key)
	}
	c.mu.Unlock()
	if ok {
		select {
		case waiter <- *ack:
		default:
		}
	}
	return &cameraAckResponse{Code: 0, Message: "ok"}, nil
}

func (c *CameraGrpc) PublishStart(node string, message subscribeMessage, timeout time.Duration) (*tokenAck, error) {
	if message.DeviceId == "" {
		return nil, errors.New("device_id is required")
	}
	node = strings.TrimSpace(node)
	if node == "" {
		node = "default"
	}
	sessionKey := c.sessionKey(message.DeviceId, message.ChannelId)
	waiter := make(chan tokenAck, 1)
	c.mu.Lock()
	c.tokenWaiters[sessionKey] = waiter
	c.mu.Unlock()

	if err := c.publish(node, message); err != nil {
		c.mu.Lock()
		delete(c.tokenWaiters, sessionKey)
		c.mu.Unlock()
		return nil, err
	}

	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	select {
	case ack := <-waiter:
		return &ack, nil
	case <-time.After(timeout):
		c.mu.Lock()
		delete(c.tokenWaiters, sessionKey)
		c.mu.Unlock()
		return nil, errors.New("wait token timeout")
	}
}

func (c *CameraGrpc) GetLastToken(deviceID, channelID string) (*tokenAck, bool) {
	key := c.sessionKey(deviceID, channelID)
	c.mu.RLock()
	ack, ok := c.lastAck[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	return &ack, true
}

func (c *CameraGrpc) addSubscriber(node string, subID uint64, ch chan subscribeMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	nodeSubscribers, ok := c.subscribers[node]
	if !ok {
		nodeSubscribers = make(map[uint64]chan subscribeMessage)
		c.subscribers[node] = nodeSubscribers
	}
	nodeSubscribers[subID] = ch
}

func (c *CameraGrpc) removeSubscriber(node string, subID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	nodeSubscribers, ok := c.subscribers[node]
	if !ok {
		return
	}
	delete(nodeSubscribers, subID)
	if len(nodeSubscribers) == 0 {
		delete(c.subscribers, node)
	}
}

func (c *CameraGrpc) publish(node string, message subscribeMessage) error {
	c.mu.RLock()
	nodeSubscribers, ok := c.subscribers[node]
	if !ok || len(nodeSubscribers) == 0 {
		c.mu.RUnlock()
		return errors.New("camera subscriber not found")
	}
	channels := make([]chan subscribeMessage, 0, len(nodeSubscribers))
	for _, ch := range nodeSubscribers {
		channels = append(channels, ch)
	}
	c.mu.RUnlock()
	for _, ch := range channels {
		select {
		case ch <- message:
		default:
		}
	}
	return nil
}

func (c *CameraGrpc) sessionKey(deviceID, channelID string) string {
	return deviceID + ":" + channelID
}

func _WebRTCService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(subscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(*CameraGrpc).Subscribe(m, &webRTCSubscribeServerWrapper{stream})
}

func _WebRTCService_SendToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(tokenAck)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*CameraGrpc).SendToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gateway.camera.WebRTCService/SendToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*CameraGrpc).SendToken(ctx, req.(*tokenAck))
	}
	return interceptor(ctx, in, info, handler)
}

var _WebRTCService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gateway.camera.WebRTCService",
	HandlerType: (*webRTCServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendToken",
			Handler:    _WebRTCService_SendToken_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _WebRTCService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "camera_grpc",
}
