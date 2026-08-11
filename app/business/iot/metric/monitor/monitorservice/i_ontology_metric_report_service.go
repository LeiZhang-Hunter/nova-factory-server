// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorservice

import (
	"context"

	modelrequest "nova-factory-server/app/business/iot/metric/monitor/monitormodels/request"
)

// IOntologyMetricReportService handles metric reports posted by ontology collectors.
type IOntologyMetricReportService interface {
	ReceiveOntologyReports(ctx context.Context, reports []modelrequest.OntologyMetricReport) (int, error)
	ReceiveExportPayload(ctx context.Context, payload modelrequest.ExportDML) error
}
