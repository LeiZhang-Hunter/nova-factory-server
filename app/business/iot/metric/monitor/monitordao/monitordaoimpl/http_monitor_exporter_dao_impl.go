// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package monitordaoimpl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	modelrequest "nova-factory-server/app/business/iot/metric/monitor/monitormodels/request"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
)

// HTTPMonitorExporterDao writes monitor payloads through HTTP.
type HTTPMonitorExporterDao struct {
	clientConfig      confighttp.ClientConfig
	logger            *zap.Logger
	telemetrySettings component.TelemetrySettings
	client            *http.Client
}

func NewHTTPMonitorExporterDao(clientConfig confighttp.ClientConfig, logger *zap.Logger, telemetrySettings component.TelemetrySettings) *HTTPMonitorExporterDao {
	return &HTTPMonitorExporterDao{
		clientConfig:      clientConfig,
		logger:            logger,
		telemetrySettings: telemetrySettings,
	}
}

func (d *HTTPMonitorExporterDao) Start(ctx context.Context, host component.Host) error {
	if d.clientConfig.Endpoint == "" {
		return fmt.Errorf("monitor exporter endpoint is required")
	}

	var extensions map[component.ID]component.Component
	if host != nil {
		extensions = host.GetExtensions()
	}

	client, err := d.clientConfig.ToClient(ctx, extensions, d.telemetrySettings)
	if err != nil {
		return fmt.Errorf("create monitor exporter http client: %w", err)
	}
	d.client = client
	return nil
}

func (d *HTTPMonitorExporterDao) Shutdown(context.Context) error {
	if d.client != nil {
		d.client.CloseIdleConnections()
	}
	return nil
}

func (d *HTTPMonitorExporterDao) Export(ctx context.Context, data modelrequest.ExportDML) error {
	if d.client == nil {
		return fmt.Errorf("monitor exporter http client is not started")
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal monitor export payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.clientConfig.Endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create monitor export request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for name, value := range d.clientConfig.Headers.Iter {
		req.Header.Set(name, string(value))
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("send monitor export request: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			d.logger.Debug("close monitor export response body failed", zap.Error(closeErr))
		}
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("monitor export request failed with status %s", resp.Status)
	}
	return nil
}
