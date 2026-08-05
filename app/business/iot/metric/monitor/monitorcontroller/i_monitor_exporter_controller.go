// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorcontroller

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// IMonitorExporterController handles Collector consume calls.
type IMonitorExporterController interface {
	ConsumeTraces(ctx context.Context, td ptrace.Traces) error
	ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error
	ConsumeLogs(ctx context.Context, ld plog.Logs) error
}
