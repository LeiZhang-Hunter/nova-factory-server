// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorserviceimpl

import (
	"context"

	"go.uber.org/zap"

	modelrequest "nova-factory-server/app/business/iot/metric/monitor/monitormodels/request"
)

type OntologyMetricReportServiceImpl struct {
	logger *zap.Logger
}

func NewOntologyMetricReportService(logger *zap.Logger) *OntologyMetricReportServiceImpl {
	return &OntologyMetricReportServiceImpl{logger: logger}
}

func (s *OntologyMetricReportServiceImpl) ReceiveOntologyReports(_ context.Context, reports []modelrequest.OntologyMetricReport) (int, error) {
	for _, report := range reports {
		s.logger.Info(
			"ontology metric report received",
			zap.String("sink", report.Sink),
			zap.String("pipeline", report.Pipeline),
			zap.String("entity_id", report.EntityID),
			zap.String("entity_type", report.EntityType),
			zap.String("file_name", report.FileName),
			zap.Int("metric_count", len(report.Metrics)),
			zap.Int("summary_count", len(report.Summary)),
			zap.Int("item_count", len(report.Items)),
		)
	}
	return len(reports), nil
}

func (s *OntologyMetricReportServiceImpl) ReceiveExportPayload(_ context.Context, payload modelrequest.ExportDML) error {
	s.logger.Info(
		"monitor export payload received",
		zap.String("signal", payload.Signal),
		zap.String("dataset", payload.Dataset),
		zap.String("format", payload.Format),
		zap.Int("body_bytes", len(payload.Body)),
	)
	return nil
}
