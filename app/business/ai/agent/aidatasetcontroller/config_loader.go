package aidatasetcontroller

import (
	"context"
	"encoding/json"
	"errors"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"nova-factory-server/app/constant/aiagent"
	context2 "nova-factory-server/app/utils/grpc/context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/gateway/gatewayservice"
	"nova-factory-server/app/constant/agent"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	v1 "nova-factory-server/app/utils/grpc/confighotload/v1"
	"nova-factory-server/app/utils/uuid"
	"strconv"
	"strings"
	"sync"
)

var configLoaderRunOnce sync.Once

type configLoaderCodec struct{}

func (configLoaderCodec) Name() string {
	return "json"
}

func (configLoaderCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (configLoaderCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type ConfigLoader struct {
	manager             *AgentConfigManager
	registry            *AgentRegistryManager
	service             gatewayservice.IAIAgentOrchestrationService
	gatewayService      gatewayservice.IAIGatewayService
	configLoaderService aidatasetservice.IConfigLoaderService
	agentService        gatewayservice.IAIAgentService

	v1.UnimplementedAgentControllerServiceServer
}

func NewConfigLoaderGrpc(service gatewayservice.IAIAgentOrchestrationService,
	gatewayService gatewayservice.IAIGatewayService,
	configLoaderService aidatasetservice.IConfigLoaderService,
	agentService gatewayservice.IAIAgentService) *ConfigLoader {
	return &ConfigLoader{
		manager:             NewAgentConfigManager(),
		registry:            NewAgentRegistryManager(),
		service:             service,
		agentService:        agentService,
		gatewayService:      gatewayService,
		configLoaderService: configLoaderService,
	}
}

// NewAgentConfigPublisher 提供智能体配置发布器实现。
func NewAgentConfigPublisher(loader *ConfigLoader) gatewayservice.IAIAgentConfigPublisher {
	return loader
}

// PrivateRoutes 注册 CameraService gRPC 服务。
func (c *ConfigLoader) PrivateRoutes(router *grpcx.GrpcServer) {
	configLoaderRunOnce.Do(func() {
		encoding.RegisterCodec(configLoaderCodec{})
	})
	v1.RegisterAgentControllerServiceServer(router.Server, c)
	return
}

func (c *ConfigLoader) PrivateHttpRoutes(router *gin.RouterGroup) {
	ai := router.Group("/ai/agent/config")
	ai.POST("/publish", middlewares.HasPermission("ai:agent:config:publish"), c.Publish)
}

// AgentReg 注册请求
func (c *ConfigLoader) AgentReg(ctx context.Context, req *v1.AgentRegQuest) (*v1.AgentRegResponse, error) {
	gatewayId, err := context2.GetGatewayId(ctx)
	if err != nil {
		zap.L().Debug("AgentHeartbeat failed", zap.Error(err))
		return nil, err
	}
	if c.registry.IsRegistered(gatewayId) {
		return &v1.AgentRegResponse{Code: 0, Msg: "already registered"}, nil
	}

	username, err := context2.GetUsername(ctx)
	if err != nil {
		return nil, errors.New("username not exist")
	}

	password, err := context2.GetPassword(ctx)
	if err != nil {
		return nil, errors.New("password not exist")
	}

	info, err := c.gatewayService.GetByID(&gin.Context{}, gatewayId)
	if err != nil {
		return &v1.AgentRegResponse{Code: -1, Msg: err.Error()}, nil
	}

	if info == nil {
		return &v1.AgentRegResponse{Code: -1, Msg: "gateway id does not exist"}, nil
	}

	if info.Username != username {
		return &v1.AgentRegResponse{Code: -1, Msg: "username not match"}, nil
	}

	if info.Password != password {
		return &v1.AgentRegResponse{Code: -1, Msg: "password not match"}, nil
	}

	c.registry.Register(gatewayId, username)

	return &v1.AgentRegResponse{Code: 0, Msg: "ok"}, nil
}

// AgentHeartbeat agent 心跳
func (c *ConfigLoader) AgentHeartbeat(ctx context.Context, req *v1.AgentHeartbeatReq) (*v1.AgentHeartbeatRes, error) {
	id, err := context2.GetGatewayId(ctx)
	if err != nil {
		zap.L().Debug("AgentHeartbeat failed", zap.Error(err))
		return nil, err
	}
	if !c.registry.IsRegistered(id) {
		return nil, errors.New("gateway id does not exist")
	}

	info, err := c.agentService.GetConfigVersion(ctx, int64(req.AgentId))
	if err != nil {
		return nil, err
	}
	if info == nil {
		return &v1.AgentHeartbeatRes{
			AgentId: req.AgentId,
			Version: "",
		}, nil
	}
	return &v1.AgentHeartbeatRes{
		AgentId: req.AgentId,
		Version: info.ConfigVersion,
	}, nil
}

// AgentGetConfig 获取配置
func (c *ConfigLoader) AgentGetConfig(ctx context.Context, req *v1.AgentGetConfigReq) (*v1.AgentGetConfigRes, error) {
	gatewayId, err := context2.GetGatewayId(ctx)
	if err != nil {
		zap.L().Debug("AgentHeartbeat failed", zap.Error(err))
		return nil, err
	}
	if !c.registry.IsRegistered(gatewayId) {
		return nil, errors.New("gateway id does not exist")
	}
	version := strings.TrimSpace(req.GetVersion())
	if version == "" {
		return nil, errors.New("config uuid does not exist")
	}

	info, err := c.configLoaderService.GetByAgentIdAndVersion(ctx, req.AgentId, version)
	if err != nil {
		zap.L().Error("AgentGetConfig failed", zap.Error(err))
		return nil, err
	}
	if info == nil {
		return nil, errors.New("config does not exist")
	}

	return &v1.AgentGetConfigRes{
		Content: info.ConfigSnapshot,
	}, nil
}

// WatchAgentChanges 双向流订阅变更
func (c *ConfigLoader) WatchAgentChanges(request *v1.WatchAgentRequest, stream v1.AgentControllerService_WatchAgentChangesServer) error {
	gatewayId := strconv.FormatInt(request.GatewayId, 10)
	c.manager.AddClient(gatewayId, stream)
	defer c.manager.DeleteClient(gatewayId, stream)

	// 读取所有的agent，下发配置
	all, err := c.configLoaderService.All(context.Background())
	if err != nil {
		zap.L().Error("WatchAgentChanges failed", zap.Error(err))
	}

	if err == nil && len(all) != 0 {
		for _, config := range all {
			message := v1.ConfigChangeEvent{
				Ref: &v1.AgentRef{
					Id:   strconv.FormatInt(config.AgentID, 10),
					Type: aiagent.AgentType,
				},
				NewVersion: config.Version,
				ConfigHash: config.ConfigMd5,
				Action:     string(aiagent.ConfigInitType),
				Content:    config.ConfigSnapshot,
				Timestamp:  timestamppb.Now(),
			}
			err := stream.Send(&message)
			if err != nil {
				zap.L().Error("send stream error", zap.Error(err))
				continue
			}
		}

	}
	//c.orchestrationService.All()

	<-stream.Context().Done()
	c.registry.Unregister(request.GatewayId)
	return nil
}

// AgentBroadcast 接收广播请求并转发到下游节点。
func (c *ConfigLoader) AgentBroadcast(ctx context.Context, req *v1.AgentBroadcastRequest) (*v1.AgentBroadcastReply, error) {
	if req == nil {
		return &v1.AgentBroadcastReply{Code: -1, Msg: "request is nil"}, nil
	}
	if req.AgentId == "" {
		return &v1.AgentBroadcastReply{Code: -1, Msg: "AgentId is empty"}, nil
	}

	id, err := strconv.ParseInt(req.AgentId, 10, 64)
	if err != nil {
		zap.L().Error("parse agent id", zap.Error(err))
		return &v1.AgentBroadcastReply{Code: -1, Msg: err.Error()}, nil
	}

	gatewayInfo, err := c.gatewayService.GetByID(&gin.Context{}, req.GatewayId)
	if err != nil {
		return nil, err
	}

	if gatewayInfo == nil {
		return nil, errors.New("gateway id does not exist")
	}

	if gatewayInfo.Username != req.Username {
		return nil, errors.New("username not match")
	}

	if gatewayInfo.Password != req.Password {
		return nil, errors.New("password not match")
	}

	info, err := c.service.GetConfigInfo(ctx, id)
	if err != nil {
		zap.L().Error("get config info", zap.Error(err))
		return &v1.AgentBroadcastReply{Code: -1, Msg: err.Error()}, nil
	}

	message := v1.ConfigChangeEvent{
		Ref: &v1.AgentRef{
			Id:   req.AgentId,
			Type: aiagent.AgentType,
		},
		NewVersion: info.ConfigMd5,
		ConfigHash: info.ConfigMd5,
		Action:     req.Action,
		Content:    info.Config,
		Timestamp:  timestamppb.Now(),
	}

	_, err = c.manager.PublishToNode(strconv.FormatInt(req.GetGatewayId(), 10), &message)
	if err != nil {
		return &v1.AgentBroadcastReply{Code: -1, Msg: err.Error()}, nil
	}

	return &v1.AgentBroadcastReply{
		Code: 0,
		Msg:  "ok",
	}, nil
}

// Publish 发布智能体配置变更
// @Summary 发布智能体配置变更
// @Description 向集群节点广播智能体配置变更通知
// @Tags 工业智能体/智能体配置
// @Param object body aidatasetmodels.AgentPublishReq true "智能体发布参数"
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.ResponseData "发布成功"
// @Router /ai/agent/config/publish [post]
func (c *ConfigLoader) Publish(ctx *gin.Context) {
	req := new(aidatasetmodels.AgentPublishReq)
	if err := ctx.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(ctx)
		return
	}

	id, err := strconv.ParseInt(strings.TrimSpace(req.AgentId), 10, 64)
	if err != nil || id == 0 {
		baizeContext.Waring(ctx, "agent id invalid")
		return
	}

	broadcastReq := &v1.AgentBroadcastRequest{
		RequestId: uuid.MakeUuid(),
		GatewayId: req.GatewayId,
		AgentId:   req.AgentId,
		Action:    req.Action,
	}

	reply, err := c.BroadcastByGrpcClient(ctx, broadcastReq)
	if err != nil {
		baizeContext.Waring(ctx, err.Error())
		return
	}
	if reply == nil || reply.GetCode() != 0 {
		baizeContext.Waring(ctx, reply.GetMsg())
		return
	}

	baizeContext.SuccessData(ctx, &aidatasetmodels.AgentPublishRes{
		Code:    reply.GetCode(),
		Message: reply.GetMsg(),
	})
}

// BroadcastByGrpcClient 通过集群地址列表并发调用广播接口。
func (c *ConfigLoader) BroadcastByGrpcClient(ctx *gin.Context, req *v1.AgentBroadcastRequest) (*v1.AgentBroadcastReply, error) {
	addressList := viper.GetStringSlice("daemonize.server_list")
	if len(addressList) == 0 {
		return nil, errors.New("daemonize.server_list is empty")
	}

	id, err := strconv.ParseInt(strings.TrimSpace(req.GetAgentId()), 10, 64)
	if err != nil || id == 0 {
		return nil, errors.New("agent id invalid")
	}
	info, err := c.service.GetConfigInfo(context.Background(), id)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, errors.New("agent not found")
	}

	gatewayInfo, err := c.gatewayService.GetByID(ctx, req.GatewayId)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(addressList))
	var mu sync.Mutex
	var replayErr error
	var successReply *v1.AgentBroadcastReply
	for _, address := range addressList {
		targetAddress := address
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			grpcCtx := metadata.AppendToOutgoingContext(ctx,
				agent.USERNAME, gatewayInfo.Username,
				agent.PASSWORD, gatewayInfo.Password,
				agent.GATEWAYID, strconv.FormatInt(gatewayInfo.ID, 10),
			)
			var reply *v1.AgentBroadcastReply
			rpcClient := v1.NewAgentControllerServiceClient(grpcx.Client.MustNewGrpcClientConn(targetAddress))
			reply, replayErr = rpcClient.AgentBroadcast(grpcCtx, req)
			if replayErr != nil {
				zap.L().Error("广播失败", zap.Error(err))
				return
			}
			if reply == nil || reply.GetCode() != 0 {
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
	return nil, errors.New("broadcast to agent subscribers failed")
}
