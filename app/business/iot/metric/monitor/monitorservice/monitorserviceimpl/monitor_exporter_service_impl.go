// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorserviceimpl

import (
	"context"
	"time"

	"nova-factory-server/app/business/iot/metric/monitor/monitordao"
	modelrequest "nova-factory-server/app/business/iot/metric/monitor/monitormodels/request"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// ConvertConfig controls OpenTelemetry signal conversion.
type ConvertConfig struct {
	Dataset string
}

type MonitorExporterServiceImpl struct {
	dao              monitordao.IMonitorExporterDao
	tracesConverter  *OtelTracesToLineProtocol
	metricsConverter *OtelMetricsToLineProtocol
	logsConverter    *OtelLogsToLineProtocol
}

func NewMonitorExporterService(dao monitordao.IMonitorExporterDao, cfg ConvertConfig) *MonitorExporterServiceImpl {
	return &MonitorExporterServiceImpl{
		dao:              dao,
		tracesConverter:  NewOtelTracesToLineProtocol(cfg),
		metricsConverter: NewOtelMetricsToLineProtocol(cfg),
		logsConverter:    NewOtelLogsToLineProtocol(cfg),
	}
}

func (s *MonitorExporterServiceImpl) ExportTraces(ctx context.Context, td ptrace.Traces) error {
	payload, err := s.tracesConverter.Convert(td)
	if err != nil {
		return err
	}
	return s.dao.Export(ctx, modelrequest.ExportDML{
		Signal:    payload.Signal,
		Dataset:   s.tracesConverter.dataset,
		Format:    payload.Format,
		CreatedAt: time.Now(),
		Body:      payload.Body,
	})
}

func (s *MonitorExporterServiceImpl) ExportMetrics(ctx context.Context, md pmetric.Metrics) error {
	payload, err := s.metricsConverter.Convert(md)
	if err != nil {
		return err
	}
	return s.dao.Export(ctx, modelrequest.ExportDML{
		Signal:    payload.Signal,
		Dataset:   s.metricsConverter.dataset,
		Format:    payload.Format,
		CreatedAt: time.Now(),
		Body:      payload.Body,
	})
}

func (s *MonitorExporterServiceImpl) ExportLogs(ctx context.Context, ld plog.Logs) error {
	payload, err := s.logsConverter.Convert(ld)
	if err != nil {
		return err
	}
	return s.dao.Export(ctx, modelrequest.ExportDML{
		Signal:    payload.Signal,
		Dataset:   s.logsConverter.dataset,
		Format:    payload.Format,
		CreatedAt: time.Now(),
		Body:      payload.Body,
	})
}
