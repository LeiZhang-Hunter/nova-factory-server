// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitordao

import (
	"context"

	modelrequest "nova-factory-server/app/business/iot/metric/monitor/monitormodels/request"
)

// IMonitorExporterDao writes converted telemetry to the monitor backend.
type IMonitorExporterDao interface {
	Export(ctx context.Context, data modelrequest.ExportDML) error
}
