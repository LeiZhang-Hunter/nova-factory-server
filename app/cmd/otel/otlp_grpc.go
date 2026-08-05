package main

import (
	"context"
	"fmt"

	monitorexporter "nova-factory-server/app/business/iot/metric/monitor"

	"github.com/spf13/viper"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.opentelemetry.io/collector/pdata/pmetric/pmetricotlp"
	"go.opentelemetry.io/collector/pdata/ptrace/ptraceotlp"
	metricnoop "go.opentelemetry.io/otel/metric/noop"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
	otlplogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	otlpmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	otlptrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type emptyComponentHost struct{}

func (emptyComponentHost) GetExtensions() map[component.ID]component.Component {
	return nil
}

type monitorOTLPBridge struct {
	traces  exporter.Traces
	metrics exporter.Metrics
	logs    exporter.Logs
}

type monitorTraceServiceServer struct {
	otlptrace.UnimplementedTraceServiceServer
	bridge *monitorOTLPBridge
}

type monitorMetricsServiceServer struct {
	otlpmetrics.UnimplementedMetricsServiceServer
	bridge *monitorOTLPBridge
}

type monitorLogsServiceServer struct {
	otlplogs.UnimplementedLogsServiceServer
	bridge *monitorOTLPBridge
}

func registerMonitorOTLPGRPC(ctx context.Context, server *grpc.Server) (func(), error) {
	factory := monitorexporter.NewFactory()
	cfg, ok := factory.CreateDefaultConfig().(*monitorexporter.Config)
	if !ok {
		return nil, fmt.Errorf("invalid monitor exporter default config type: %T", factory.CreateDefaultConfig())
	}
	applyMonitorExporterConfig(cfg)

	settings := exporter.Settings{
		ID: component.NewID(factory.Type()),
		TelemetrySettings: component.TelemetrySettings{
			Logger:         zap.L(),
			MeterProvider:  metricnoop.NewMeterProvider(),
			TracerProvider: tracenoop.NewTracerProvider(),
		},
	}

	tracesExporter, err := factory.CreateTraces(ctx, settings, cfg)
	if err != nil {
		return nil, fmt.Errorf("create monitor traces exporter: %w", err)
	}
	metricsExporter, err := factory.CreateMetrics(ctx, settings, cfg)
	if err != nil {
		return nil, fmt.Errorf("create monitor metrics exporter: %w", err)
	}
	logsExporter, err := factory.CreateLogs(ctx, settings, cfg)
	if err != nil {
		return nil, fmt.Errorf("create monitor logs exporter: %w", err)
	}

	host := emptyComponentHost{}
	if err := tracesExporter.Start(ctx, host); err != nil {
		return nil, fmt.Errorf("start monitor traces exporter: %w", err)
	}
	if err := metricsExporter.Start(ctx, host); err != nil {
		_ = tracesExporter.Shutdown(ctx)
		return nil, fmt.Errorf("start monitor metrics exporter: %w", err)
	}
	if err := logsExporter.Start(ctx, host); err != nil {
		_ = metricsExporter.Shutdown(ctx)
		_ = tracesExporter.Shutdown(ctx)
		return nil, fmt.Errorf("start monitor logs exporter: %w", err)
	}

	bridge := &monitorOTLPBridge{
		traces:  tracesExporter,
		metrics: metricsExporter,
		logs:    logsExporter,
	}
	otlptrace.RegisterTraceServiceServer(server, &monitorTraceServiceServer{bridge: bridge})
	otlpmetrics.RegisterMetricsServiceServer(server, &monitorMetricsServiceServer{bridge: bridge})
	otlplogs.RegisterLogsServiceServer(server, &monitorLogsServiceServer{bridge: bridge})

	return func() {
		shutdownCtx := context.Background()
		if err := logsExporter.Shutdown(shutdownCtx); err != nil {
			zap.L().Error("shutdown monitor logs exporter failed", zap.Error(err))
		}
		if err := metricsExporter.Shutdown(shutdownCtx); err != nil {
			zap.L().Error("shutdown monitor metrics exporter failed", zap.Error(err))
		}
		if err := tracesExporter.Shutdown(shutdownCtx); err != nil {
			zap.L().Error("shutdown monitor traces exporter failed", zap.Error(err))
		}
	}, nil
}

func applyMonitorExporterConfig(cfg *monitorexporter.Config) {
	if endpoint := viper.GetString("otel.exporter.endpoint"); endpoint != "" {
		cfg.ClientConfig.Endpoint = endpoint
	}
	if endpoint := viper.GetString("metric.monitor_exporter.endpoint"); endpoint != "" {
		cfg.ClientConfig.Endpoint = endpoint
	}
	if dataset := viper.GetString("otel.exporter.dataset"); dataset != "" {
		cfg.Dataset = dataset
	}
	if dataset := viper.GetString("metric.monitor_exporter.dataset"); dataset != "" {
		cfg.Dataset = dataset
	}
}

func (s *monitorTraceServiceServer) Export(ctx context.Context, req *otlptrace.ExportTraceServiceRequest) (*otlptrace.ExportTraceServiceResponse, error) {
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal otlp trace request: %w", err)
	}

	exportRequest := ptraceotlp.NewExportRequest()
	if err := exportRequest.UnmarshalProto(data); err != nil {
		return nil, fmt.Errorf("unmarshal otlp trace request: %w", err)
	}
	if err := s.bridge.traces.ConsumeTraces(ctx, exportRequest.Traces()); err != nil {
		return nil, err
	}
	return &otlptrace.ExportTraceServiceResponse{}, nil
}

func (s *monitorMetricsServiceServer) Export(ctx context.Context, req *otlpmetrics.ExportMetricsServiceRequest) (*otlpmetrics.ExportMetricsServiceResponse, error) {
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal otlp metrics request: %w", err)
	}

	exportRequest := pmetricotlp.NewExportRequest()
	if err := exportRequest.UnmarshalProto(data); err != nil {
		return nil, fmt.Errorf("unmarshal otlp metrics request: %w", err)
	}
	if err := s.bridge.metrics.ConsumeMetrics(ctx, exportRequest.Metrics()); err != nil {
		return nil, err
	}
	return &otlpmetrics.ExportMetricsServiceResponse{}, nil
}

func (s *monitorLogsServiceServer) Export(ctx context.Context, req *otlplogs.ExportLogsServiceRequest) (*otlplogs.ExportLogsServiceResponse, error) {
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal otlp logs request: %w", err)
	}

	exportRequest := plogotlp.NewExportRequest()
	if err := exportRequest.UnmarshalProto(data); err != nil {
		return nil, fmt.Errorf("unmarshal otlp logs request: %w", err)
	}
	if err := s.bridge.logs.ConsumeLogs(ctx, exportRequest.Logs()); err != nil {
		return nil, err
	}
	return &otlplogs.ExportLogsServiceResponse{}, nil
}
