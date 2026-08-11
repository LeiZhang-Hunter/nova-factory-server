// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitorserviceimpl

import (
	modelentity "nova-factory-server/app/business/iot/metric/monitor/monitormodels/entity"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

type OtelMetricsToLineProtocol struct {
	dataset string
}

func NewOtelMetricsToLineProtocol(cfg ConvertConfig) *OtelMetricsToLineProtocol {
	return &OtelMetricsToLineProtocol{dataset: cfg.Dataset}
}

func (c *OtelMetricsToLineProtocol) Convert(md pmetric.Metrics) (modelentity.ExportPayload, error) {
	body, err := (&pmetric.JSONMarshaler{}).MarshalMetrics(md)
	if err != nil {
		return modelentity.ExportPayload{}, err
	}
	return modelentity.ExportPayload{
		Signal: "metrics",
		Format: "otlp_json",
		Body:   body,
	}, nil
}
