package daemonizecontroller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	xdsv1 "github.com/novawatcher-io/nova-factory-payload/gateway/xds/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const pipelineSnapshotTypeURL = "gateway.xds.v1.PipelineSnapshot"

// PrivateGrpcRoutes registers the dedicated gateway xDS service. It is kept
// separate from the old daemonize stream so dynamic runtime config can evolve
// independently from process-management traffic.
func (c *Config) PrivateGrpcRoutes(router *grpcx.GrpcServer) {
	xdsv1.RegisterGatewayDiscoveryServiceServer(router.Server, c)
}

// FetchResources returns the latest immutable pipeline snapshot for a gateway.
// The first implementation focuses on PipelineSnapshot because the backend
// compiler already renders source, sink and alert rules into one runtime graph.
func (c *Config) FetchResources(ctx context.Context, req *xdsv1.DiscoveryRequest) (*xdsv1.DiscoveryResponse, error) {
	gatewayID, err := c.authorizeGateway(ctx, req.GetGatewayId())
	if err != nil {
		return nil, err
	}
	return c.buildDiscoveryResponse(ctx, gatewayID)
}

// StreamResources keeps a long-lived control-plane stream. The server polls the
// current compiled snapshot periodically; when the rendered hash changes it
// pushes a new response to the gateway.
func (c *Config) StreamResources(stream xdsv1.GatewayDiscoveryService_StreamResourcesServer) error {
	firstReq, err := stream.Recv()
	if err != nil {
		return err
	}

	gatewayID, err := c.authorizeGateway(stream.Context(), firstReq.GetGatewayId())
	if err != nil {
		return err
	}

	if err = c.ensureTypeURL(firstReq.GetTypeUrl()); err != nil {
		return err
	}

	ackVersion := firstReq.GetVersionInfo()
	sentVersion := ""
	lastProbeVersion := ""

	reqCh := make(chan *xdsv1.DiscoveryRequest, 4)
	errCh := make(chan error, 1)
	go func() {
		for {
			req, recvErr := stream.Recv()
			if recvErr != nil {
				errCh <- recvErr
				return
			}
			reqCh <- req
		}
	}()

	sendLatest := func() error {
		probeVersion, probeErr := c.getLatestPipelineProbeVersion(stream.Context(), gatewayID)
		if probeErr != nil {
			return probeErr
		}
		if probeVersion == "" || probeVersion == lastProbeVersion {
			return nil
		}

		resp, buildErr := c.buildDiscoveryResponse(stream.Context(), gatewayID)
		if buildErr != nil {
			return buildErr
		}
		if resp.GetVersionInfo() == "" {
			return nil
		}
		if resp.GetVersionInfo() == ackVersion || resp.GetVersionInfo() == sentVersion {
			lastProbeVersion = probeVersion
			return nil
		}
		if sendErr := stream.Send(resp); sendErr != nil {
			return sendErr
		}
		lastProbeVersion = probeVersion
		sentVersion = resp.GetVersionInfo()
		return nil
	}

	if err = sendLatest(); err != nil {
		return err
	}

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case recvErr := <-errCh:
			if errors.Is(recvErr, io.EOF) {
				return nil
			}
			return recvErr
		case req := <-reqCh:
			if req == nil {
				continue
			}
			if err = c.ensureTypeURL(req.GetTypeUrl()); err != nil {
				return err
			}
			if req.GetAckStatus() == "ACK" {
				ackVersion = req.GetVersionInfo()
				continue
			}
			if req.GetAckStatus() == "NACK" {
				zap.L().Warn("gateway xds nack", zap.String("gatewayId", strconv.FormatUint(gatewayID, 10)), zap.String("version", req.GetVersionInfo()), zap.String("message", req.GetErrorMessage()))
			}
		case <-ticker.C:
			if err = sendLatest(); err != nil {
				return err
			}
		}
	}
}

func (c *Config) buildDiscoveryResponse(ctx context.Context, gatewayID uint64) (*xdsv1.DiscoveryResponse, error) {
	snapshot, err := c.buildPipelineSnapshot(ctx, gatewayID)
	if err != nil {
		return nil, err
	}

	body, err := proto.Marshal(snapshot)
	if err != nil {
		return nil, err
	}

	nonce := fmt.Sprintf("%s-%d", snapshot.GetVersionInfo(), time.Now().UnixNano())
	resource := &xdsv1.Resource{
		Name:            fmt.Sprintf("gateway-%d-pipelines", gatewayID),
		ResourceVersion: snapshot.GetVersionInfo(),
		TypeUrl:         pipelineSnapshotTypeURL,
		Body:            body,
	}

	return &xdsv1.DiscoveryResponse{
		VersionInfo: snapshot.GetVersionInfo(),
		TypeUrl:     pipelineSnapshotTypeURL,
		Nonce:       nonce,
		Resources:   []*xdsv1.Resource{resource},
	}, nil
}

func (c *Config) buildPipelineSnapshot(ctx context.Context, gatewayID uint64) (*xdsv1.PipelineSnapshot, error) {
	config, err := c.configService.GetLastedConfig(ctx, gatewayID)
	if err != nil {
		return nil, err
	}

	return &xdsv1.PipelineSnapshot{
		TypeUrl:     pipelineSnapshotTypeURL,
		GatewayId:   strconv.FormatUint(gatewayID, 10),
		VersionInfo: config.ConfigVersion,
		YamlContent: config.Content,
		ContentHash: config.ContentHash,
		GeneratedAt: timestamppb.Now(),
	}, nil
}

func (c *Config) getLatestPipelineProbeVersion(ctx context.Context, gatewayID uint64) (string, error) {
	config, err := c.configService.GetLastedConfigHashAndVersion(ctx, gatewayID)
	if err != nil {
		return "", err
	}
	if config == nil {
		return "", nil
	}
	if config.ContentHash != "" {
		return config.ContentHash, nil
	}
	return config.ConfigVersion, nil
}

func (c *Config) authorizeGateway(ctx context.Context, gatewayID string) (uint64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("grpc metadata missing")
	}

	username := firstValue(md.Get("username"))
	password := firstValue(md.Get("password"))
	if gatewayID == "" {
		gatewayID = firstValue(md.Get("gateway_id"))
	}
	if gatewayID == "" {
		gatewayID = firstValue(md.Get("gatewayid"))
	}
	if username == "" || password == "" || gatewayID == "" {
		return 0, errors.New("gateway credentials missing")
	}

	parsedGatewayID, err := strconv.ParseUint(gatewayID, 10, 64)
	if err != nil {
		return 0, err
	}

	info, err := c.agentService.GetByObjectId(ctx, parsedGatewayID)
	if err != nil {
		return 0, err
	}
	if info == nil {
		return 0, errors.New("gateway not found")
	}
	if info.Username != username || info.Password != password {
		return 0, errors.New("gateway credential validation failed")
	}

	return parsedGatewayID, nil
}

func (c *Config) ensureTypeURL(typeURL string) error {
	if typeURL == "" || typeURL == pipelineSnapshotTypeURL {
		return nil
	}
	return fmt.Errorf("unsupported type url: %s", typeURL)
}

func firstValue(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
