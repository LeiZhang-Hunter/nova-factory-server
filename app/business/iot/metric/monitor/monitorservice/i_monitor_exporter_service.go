// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorservice

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// IMonitorExporterService converts and exports OpenTelemetry signals.
type IMonitorExporterService interface {
	ExportTraces(ctx context.Context, td ptrace.Traces) error
	ExportMetrics(ctx context.Context, md pmetric.Metrics) error
	ExportLogs(ctx context.Context, ld plog.Logs) error
}
