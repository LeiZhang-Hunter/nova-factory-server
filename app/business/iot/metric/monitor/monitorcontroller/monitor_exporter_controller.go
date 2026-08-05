// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorcontroller

import (
	"context"

	"go.uber.org/zap"

	"nova-factory-server/app/business/iot/metric/monitor/monitorservice"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type MonitorExporterController struct {
	service monitorservice.IMonitorExporterService
	logger  *zap.Logger
}

func NewMonitorExporterController(service monitorservice.IMonitorExporterService, logger *zap.Logger) *MonitorExporterController {
	return &MonitorExporterController{
		service: service,
		logger:  logger,
	}
}

func (c *MonitorExporterController) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	if err := c.service.ExportTraces(ctx, td); err != nil {
		c.logger.Error("export traces failed", zap.Error(err))
		return err
	}
	return nil
}

func (c *MonitorExporterController) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	if err := c.service.ExportMetrics(ctx, md); err != nil {
		c.logger.Error("export metrics failed", zap.Error(err))
		return err
	}
	return nil
}

func (c *MonitorExporterController) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	if err := c.service.ExportLogs(ctx, ld); err != nil {
		c.logger.Error("export logs failed", zap.Error(err))
		return err
	}
	return nil
}
