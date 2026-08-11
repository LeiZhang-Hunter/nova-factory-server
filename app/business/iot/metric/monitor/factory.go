// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package monitorexporter implements the OpenTelemetry Collector exporter.
//
//go:generate make mdatagen
package monitorexporter // import "nova-factory-server/app/business/iot/metric/monitor"

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"nova-factory-server/app/business/iot/metric/monitor/internal/metadata"
	"nova-factory-server/app/business/iot/metric/monitor/monitorcontroller"
	"nova-factory-server/app/business/iot/metric/monitor/monitordao/monitordaoimpl"
	"nova-factory-server/app/business/iot/metric/monitor/monitorservice/monitorserviceimpl"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configoptional"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const defaultEndpoint = "http://localhost:6000/otel/monitor"

// NewFactory creates a factory for the monitor exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithTraces(createTraceExporter, metadata.TracesStability),
		exporter.WithMetrics(createMetricsExporter, metadata.MetricsStability),
		exporter.WithLogs(createLogsExporter, metadata.LogsStability),
	)
}

func createDefaultConfig() component.Config {
	clientConfig := confighttp.NewDefaultClientConfig()
	clientConfig.Endpoint = defaultEndpoint
	clientConfig.Timeout = 30 * time.Second
	clientConfig.Headers = configopaque.MapList{
		{Name: "Content-Type", Value: configopaque.String("application/json")},
	}

	return &Config{
		ClientConfig:  clientConfig,
		QueueSettings: configoptional.Some(exporterhelper.NewDefaultQueueConfig()),
		BackOffConfig: configretry.NewDefaultBackOffConfig(),
		Dataset:       "default",
	}
}

func createTraceExporter(ctx context.Context, set exporter.Settings, config component.Config) (exporter.Traces, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid monitor exporter config type: %T", config)
	}

	logger := set.Logger.With(zap.String("signal", "traces"))
	writer := monitordaoimpl.NewHTTPMonitorExporterDao(cfg.ClientConfig, logger, set.TelemetrySettings)
	service := monitorserviceimpl.NewMonitorExporterService(writer, monitorserviceimpl.ConvertConfig{
		Dataset: cfg.Dataset,
	})
	controller := monitorcontroller.NewMonitorExporterController(service, logger)

	return exporterhelper.NewTraces(
		ctx,
		set,
		config,
		controller.ConsumeTraces,
		exporterhelper.WithQueue(cfg.QueueSettings),
		exporterhelper.WithRetry(cfg.BackOffConfig),
		exporterhelper.WithStart(writer.Start),
		exporterhelper.WithShutdown(writer.Shutdown),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
	)
}

func createMetricsExporter(ctx context.Context, set exporter.Settings, config component.Config) (exporter.Metrics, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid monitor exporter config type: %T", config)
	}

	logger := set.Logger.With(zap.String("signal", "metrics"))
	writer := monitordaoimpl.NewHTTPMonitorExporterDao(cfg.ClientConfig, logger, set.TelemetrySettings)
	service := monitorserviceimpl.NewMonitorExporterService(writer, monitorserviceimpl.ConvertConfig{
		Dataset: cfg.Dataset,
	})
	controller := monitorcontroller.NewMonitorExporterController(service, logger)

	return exporterhelper.NewMetrics(
		ctx,
		set,
		config,
		controller.ConsumeMetrics,
		exporterhelper.WithQueue(cfg.QueueSettings),
		exporterhelper.WithRetry(cfg.BackOffConfig),
		exporterhelper.WithStart(writer.Start),
		exporterhelper.WithShutdown(writer.Shutdown),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
	)
}

func createLogsExporter(ctx context.Context, set exporter.Settings, config component.Config) (exporter.Logs, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid monitor exporter config type: %T", config)
	}

	logger := set.Logger.With(zap.String("signal", "logs"))
	writer := monitordaoimpl.NewHTTPMonitorExporterDao(cfg.ClientConfig, logger, set.TelemetrySettings)
	service := monitorserviceimpl.NewMonitorExporterService(writer, monitorserviceimpl.ConvertConfig{
		Dataset: cfg.Dataset,
	})
	controller := monitorcontroller.NewMonitorExporterController(service, logger)

	return exporterhelper.NewLogs(
		ctx,
		set,
		config,
		controller.ConsumeLogs,
		exporterhelper.WithQueue(cfg.QueueSettings),
		exporterhelper.WithRetry(cfg.BackOffConfig),
		exporterhelper.WithStart(writer.Start),
		exporterhelper.WithShutdown(writer.Shutdown),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
	)
}
