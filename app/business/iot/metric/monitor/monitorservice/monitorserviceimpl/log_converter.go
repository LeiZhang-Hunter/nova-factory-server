// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorserviceimpl

import (
	modelentity "nova-factory-server/app/business/iot/metric/monitor/monitormodels/entity"

	"go.opentelemetry.io/collector/pdata/plog"
)

type OtelLogsToLineProtocol struct {
	dataset string
}

func NewOtelLogsToLineProtocol(cfg ConvertConfig) *OtelLogsToLineProtocol {
	return &OtelLogsToLineProtocol{dataset: cfg.Dataset}
}

func (c *OtelLogsToLineProtocol) Convert(ld plog.Logs) (modelentity.ExportPayload, error) {
	body, err := (&plog.JSONMarshaler{}).MarshalLogs(ld)
	if err != nil {
		return modelentity.ExportPayload{}, err
	}
	return modelentity.ExportPayload{
		Signal: "logs",
		Format: "otlp_json",
		Body:   body,
	}, nil
}
