// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorserviceimpl

import (
	modelentity "nova-factory-server/app/business/iot/metric/monitor/monitormodels/entity"

	"go.opentelemetry.io/collector/pdata/ptrace"
)

type OtelTracesToLineProtocol struct {
	dataset string
}

func NewOtelTracesToLineProtocol(cfg ConvertConfig) *OtelTracesToLineProtocol {
	return &OtelTracesToLineProtocol{dataset: cfg.Dataset}
}

func (c *OtelTracesToLineProtocol) Convert(td ptrace.Traces) (modelentity.ExportPayload, error) {
	body, err := (&ptrace.JSONMarshaler{}).MarshalTraces(td)
	if err != nil {
		return modelentity.ExportPayload{}, err
	}
	return modelentity.ExportPayload{
		Signal: "traces",
		Format: "otlp_json",
		Body:   body,
	}, nil
}
